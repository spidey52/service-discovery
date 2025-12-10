# Service Discovery UI

A modern web dashboard for monitoring and managing your service discovery system.

## Features

- **Real-time Service Monitoring**: View all registered services with their current status
- **Search & Filter**: Find services by name, ID, host, or filter by environment
- **Service Statistics**: Overview of total services, active services, and environments
- **Auto-refresh**: Automatically updates every 30 seconds
- **Responsive Design**: Works on desktop and mobile devices
- **Status Indicators**: Visual indicators for active/inactive services

## Usage

### Starting the Server

The UI is served automatically when you run the service discovery server:

```bash
go run main.go
```

The server will start on port 4000, and the UI will be available at:

- **Dashboard**: http://localhost:4000
- **Direct UI**: http://localhost:4000/ui/

### Features Overview

#### Service Cards

Each service is displayed in a card showing:

- Service name and unique ID
- Host and port information
- Environment, region, and version
- Current status (Active/Inactive)
- Last heartbeat timestamp
- Developer information (if available)

#### Search & Filtering

- **Search**: Filter services by name, ID, or host
- **Environment Filter**: Filter by dev, staging, or prod environments
- **Real-time**: Filters apply instantly as you type

#### Statistics Dashboard

- **Total Services**: All registered services
- **Active Services**: Services with recent heartbeats
- **Environments**: Number of different environments

#### Auto-refresh

- Services list refreshes every 30 seconds automatically
- Manual refresh button available
- Error handling with retry functionality

## API Integration

The UI communicates with the service discovery API endpoints:

- `GET /lookup` - Retrieves all registered services
- Services are filtered client-side for search functionality

## Development

### File Structure

```
ui/
├── index.html      # Main dashboard page
├── styles.css      # Modern responsive styling
├── app.js         # JavaScript functionality
└── README.md      # This documentation
```

### Customization

You can customize the UI by modifying:

1. **Colors & Styling**: Edit `styles.css`
2. **Functionality**: Modify `app.js`
3. **Layout**: Update `index.html`

### CORS Considerations

If you need to access the UI from a different domain, ensure your Go server has appropriate CORS headers configured.

## Browser Support

- Chrome 70+
- Firefox 65+
- Safari 12+
- Edge 79+

## Troubleshooting

### UI Not Loading

- Ensure the server is running on port 4000
- Check that the `ui/` directory exists in the project root
- Verify static file serving is enabled in `main.go`

### Services Not Showing

- Check MongoDB connection
- Verify services are registered and sending heartbeats
- Check browser console for API errors

### Search/Filter Not Working

- Ensure JavaScript is enabled
- Check browser console for errors
- Try refreshing the page

## Contributing

To contribute to the UI:

1. Make changes to the files in the `ui/` directory
2. Test locally with the Go server running
3. Ensure responsive design works on mobile
4. Update this README if adding new features
