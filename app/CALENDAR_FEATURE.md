# Multiple Calendar Feature

## Overview
The app now supports creating and managing multiple calendars, allowing you to track different aspects of your life separately.

## Features

### 1. Multiple Calendars
- Create unlimited custom calendars
- Each calendar has its own name and color scheme
- Track different metrics in separate calendars (e.g., "Mental Health", "Food Quality", "Exercise")

### 2. Calendar Management
- **Create**: Add new calendars with custom names and color schemes
- **Edit**: Modify calendar names and colors
- **Delete**: Remove calendars you no longer need
- **Select**: Switch between calendars using the dropdown selector

### 3. Color Customization
- Each calendar can have its own set of colors
- Customize colors in the Settings screen
- Colors are automatically applied when you switch calendars

### 4. User Interface
- **Calendar Selector**: Dropdown in the main calendar screen to switch between calendars
- **Settings Integration**: Manage all calendars from the Settings screen
- **One Calendar at a Time**: View and edit one calendar at a time for focused tracking

## How to Use

### Creating Your First Calendar
1. Open the app
2. Go to Settings
3. Tap "Manage Calendars"
4. Tap "Create New Calendar"
5. Enter a name (e.g., "Mental Health")
6. Customize the colors if desired
7. Tap "Create"

### Switching Between Calendars
1. On the main calendar screen, look for the dropdown at the top
2. Tap the dropdown to see all your calendars
3. Select the calendar you want to view/edit

### Managing Calendars
1. Go to Settings
2. Tap "Manage Calendars"
3. From here you can:
   - Create new calendars
   - Edit existing calendar names and colors
   - Delete calendars you no longer need

## Data Migration
- Existing day colors are automatically migrated to a default calendar called "My Calendar"
- Your data is preserved during the upgrade
- Legacy storage is cleaned up after successful migration

## Example Use Cases
- **Mental Health Calendar**: Track mood with colors representing different emotional states
- **Food Quality Calendar**: Track eating habits with colors for healthy/unhealthy days
- **Exercise Calendar**: Track workout consistency with colors for different activity levels
- **Sleep Quality Calendar**: Track sleep patterns with colors for sleep quality
- **Work Productivity Calendar**: Track productive vs. unproductive days

## Technical Notes
- Calendars are stored locally using SharedPreferences
- Each calendar maintains its own set of colored days
- The app supports backward compatibility with existing data
- Color schemes are customizable per calendar
