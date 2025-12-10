// Service Discovery Dashboard JavaScript

class ServiceDiscoveryDashboard {
 constructor() {
  this.baseUrl = window.location.origin || "http://localhost:4000";
  this.services = [];
  this.filteredServices = [];
  this.websocket = null;
  this.reconnectAttempts = 0;
  this.maxReconnectAttempts = 5;

  this.initializeElements();
  this.attachEventListeners();
  this.connectWebSocket();
  this.loadServices(); // Initial load as fallback
 }

 initializeElements() {
  this.searchInput = document.getElementById("searchInput");
  this.environmentFilter = document.getElementById("environmentFilter");
  this.refreshBtn = document.getElementById("refreshBtn");
  this.retryBtn = document.getElementById("retryBtn");
  this.servicesGrid = document.getElementById("servicesGrid");
  this.loading = document.getElementById("loading");
  this.error = document.getElementById("error");
  this.totalServicesEl = document.getElementById("totalServices");
  this.activeServicesEl = document.getElementById("activeServices");
  this.environmentsEl = document.getElementById("environments");
 }

 attachEventListeners() {
  this.searchInput.addEventListener("input", () => this.filterServices());
  this.environmentFilter.addEventListener("change", () => this.filterServices());
  this.refreshBtn.addEventListener("click", () => this.loadServices());
  this.retryBtn.addEventListener("click", () => this.loadServices());
 }

 async loadServices() {
  this.showLoading();
  this.hideError();

  try {
   const response = await fetch(`${this.baseUrl}/lookup`);
   if (!response.ok) {
    throw new Error(`HTTP ${response.status}: ${response.statusText}`);
   }

   this.services = await response.json();
   this.updateStats();
   this.filterServices();
   this.hideLoading();
  } catch (error) {
   console.error("Failed to load services:", error);
   this.showError();
   this.hideLoading();
  }
 }

 filterServices() {
  const searchTerm = this.searchInput.value.toLowerCase();
  const environmentFilter = this.environmentFilter.value;

  this.filteredServices = this.services.filter((service) => {
   const matchesSearch = service.serviceName.toLowerCase().includes(searchTerm) || service.id.toLowerCase().includes(searchTerm) || service.host.toLowerCase().includes(searchTerm);

   const matchesEnvironment = !environmentFilter || service.metadata.environment === environmentFilter;

   return matchesSearch && matchesEnvironment;
  });

  this.renderServices();
 }

 updateStats() {
  const totalServices = this.services.length;
  const activeServices = this.services.filter((s) => this.isServiceActive(s)).length;
  const environments = new Set(this.services.map((s) => s.metadata.environment)).size;

  this.totalServicesEl.textContent = totalServices;
  this.activeServicesEl.textContent = activeServices;
  this.environmentsEl.textContent = environments;
 }

 isServiceActive(service) {
  if (!service.lastHeartbeat) return false;

  const lastHeartbeat = new Date(service.lastHeartbeat);
  const now = new Date();
  const timeDiff = now - lastHeartbeat;
  const ttl = 30 * 1000; // 30 seconds TTL

  return timeDiff < ttl;
 }

 renderServices() {
  this.servicesGrid.innerHTML = "";

  if (this.filteredServices.length === 0) {
   this.servicesGrid.innerHTML = `
                <div class="service-card" style="grid-column: 1 / -1; text-align: center; padding: 40px;">
                    <h3 style="color: #666; margin-bottom: 10px;">No services found</h3>
                    <p style="color: #999;">Try adjusting your search or filter criteria.</p>
                </div>
            `;
   return;
  }

  this.filteredServices.forEach((service) => {
   const serviceCard = this.createServiceCard(service);
   this.servicesGrid.appendChild(serviceCard);
  });
 }

 createServiceCard(service) {
  const card = document.createElement("div");
  card.className = "service-card";

  const isActive = this.isServiceActive(service);
  const statusClass = isActive ? "status-active" : "status-inactive";
  const statusText = isActive ? "Active" : "Inactive";

  const lastHeartbeat = service.lastHeartbeat ? new Date(service.lastHeartbeat).toLocaleString() : "Never";

  card.innerHTML = `
            <div class="service-header">
                <h3>${service.serviceName}</h3>
                <div class="service-id">${service.id}</div>
            </div>
            <div class="service-body">
                <div class="service-info">
                    <div class="info-item">
                        <div class="info-label">Host</div>
                        <div class="info-value">${service.host}:${service.port}</div>
                    </div>
                    <div class="info-item">
                        <div class="info-label">Environment</div>
                        <div class="info-value">${service.metadata.environment}</div>
                    </div>
                    <div class="info-item">
                        <div class="info-label">Region</div>
                        <div class="info-value">${service.metadata.region}</div>
                    </div>
                    <div class="info-item">
                        <div class="info-label">Version</div>
                        <div class="info-value">${service.metadata.version}</div>
                    </div>
                    <div class="info-item">
                        <div class="info-label">Status</div>
                        <div class="info-value">
                            <span class="service-status ${statusClass}">${statusText}</span>
                        </div>
                    </div>
                    <div class="info-item">
                        <div class="info-label">Last Heartbeat</div>
                        <div class="info-value">${lastHeartbeat}</div>
                    </div>
                </div>
                ${
                 service.metadata.developer
                  ? `
                    <div class="info-item" style="grid-column: 1 / -1;">
                        <div class="info-label">Developer</div>
                        <div class="info-value">${service.metadata.developer}</div>
                    </div>
                `
                  : ""
                }
            </div>
        `;

  return card;
 }

 showLoading() {
  this.loading.style.display = "block";
  this.servicesGrid.style.display = "none";
 }

 hideLoading() {
  this.loading.style.display = "none";
  this.servicesGrid.style.display = "grid";
 }

 showError() {
  this.error.style.display = "block";
 }

 hideError() {
  this.error.style.display = "none";
 }

 connectWebSocket() {
  const wsUrl = this.baseUrl.replace(/^http/, "ws") + "/ws";
  console.log("Connecting to WebSocket:", wsUrl);

  try {
   this.websocket = new WebSocket(wsUrl);

   this.websocket.onopen = (event) => {
    console.log("WebSocket connected");
    this.reconnectAttempts = 0;
    this.hideError();
   };

   this.websocket.onmessage = (event) => {
    try {
     const data = JSON.parse(event.data);
     if (data.type === "services_update") {
      console.log("Received services update via WebSocket");
      this.loadServices(); // Reload services when update received
     }
    } catch (error) {
     console.error("Failed to parse WebSocket message:", error);
    }
   };

   this.websocket.onclose = (event) => {
    console.log("WebSocket disconnected, attempting reconnect...");
    this.attemptReconnect();
   };

   this.websocket.onerror = (error) => {
    console.error("WebSocket error:", error);
   };
  } catch (error) {
   console.error("Failed to create WebSocket connection:", error);
   this.attemptReconnect();
  }
 }

 attemptReconnect() {
  if (this.reconnectAttempts >= this.maxReconnectAttempts) {
   console.error("Max WebSocket reconnection attempts reached");
   this.showError();
   return;
  }

  this.reconnectAttempts++;
  const delay = Math.min(1000 * Math.pow(2, this.reconnectAttempts), 30000); // Exponential backoff

  console.log(`Attempting WebSocket reconnect ${this.reconnectAttempts}/${this.maxReconnectAttempts} in ${delay}ms`);

  setTimeout(() => {
   this.connectWebSocket();
  }, delay);
 }

 disconnectWebSocket() {
  if (this.websocket) {
   this.websocket.close();
   this.websocket = null;
  }
 }
}

// Fallback polling every 60 seconds (only if WebSocket fails)
setInterval(() => {
 if (window.dashboard && (!window.dashboard.websocket || window.dashboard.websocket.readyState !== WebSocket.OPEN)) {
  console.log("WebSocket not connected, falling back to polling");
  window.dashboard.loadServices();
 }
}, 1000 * 60);

// Initialize the dashboard when the page loads
document.addEventListener("DOMContentLoaded", () => {
 window.dashboard = new ServiceDiscoveryDashboard();
});
