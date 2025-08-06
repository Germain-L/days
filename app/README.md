# Days - Android Day Tracker App

A modern Android application built with Jetpack Compose that allows users to track their days by assigning colors to dates on a calendar.

## 🚀 Features

- **Interactive Calendar**: Tap any day to mark it with your selected color
- **Color Management**: Create, edit, and delete custom colors with meaningful names
- **Theme Support**: Light/Dark theme with system theme following
- **Data Persistence**: Local storage with export/import functionality
- **Clean UI**: Material 3 design with accessibility support
- **Error Handling**: Comprehensive error states with user-friendly messages

## 🏗️ Architecture

This app follows modern Android development practices:

- **MVVM Architecture** with ViewModels for business logic
- **Repository Pattern** for data abstraction
- **Clean Architecture** with clear separation of concerns
- **Jetpack Compose** for modern UI development
- **Coroutines & Flow** for reactive programming
- **Material 3** for consistent design

## 📱 Tech Stack

- **Language**: Kotlin
- **UI Framework**: Jetpack Compose
- **Architecture**: MVVM + Repository Pattern
- **Concurrency**: Kotlin Coroutines
- **Data Storage**: SharedPreferences with JSON serialization
- **Design System**: Material 3
- **Navigation**: Navigation Compose
- **Build System**: Gradle with Kotlin DSL

## 🗂️ Project Structure

```
src/main/java/com/germainleignel/days/
├── data/           # Data models and DTOs
├── navigation/     # Navigation logic and setup
├── storage/        # Data persistence layer
│   ├── DataRepository.kt      # Repository interface
│   ├── LocalDataRepository.kt # Local storage implementation
│   ├── StorageModels.kt       # Serialization models
│   └── RepositoryFactory.kt   # Repository factory
├── ui/
│   ├── components/ # Reusable UI components
│   ├── screens/    # Screen composables
│   └── theme/      # Theme and styling
├── viewmodel/      # ViewModels for business logic
└── MainActivity.kt # Entry point
```

## 🛠️ Building the Project

### Prerequisites
- Android Studio Hedgehog (2023.1.1) or later
- Android SDK 35 or higher
- Kotlin 1.9.20 or later

### Build Variants
- **Debug**: Development build with debugging enabled
- **Release**: Production build with ProGuard optimization

### Build Commands
```bash
# Debug build
./gradlew assembleDebug

# Release build
./gradlew assembleRelease

# Run tests
./gradlew test

# Run on device/emulator
./gradlew installDebug
```

## 📋 Key Components

### ViewModels
- `DayTrackerViewModel`: Manages app state, day colors, and settings

### Repositories
- `DataRepository`: Interface for data operations
- `LocalDataRepository`: SharedPreferences implementation
- `RepositoryFactory`: Factory for repository instances

### UI Components
- `CalendarScreen`: Main calendar interface
- `SettingsScreen`: App configuration and data management
- `DayTrackerComponents`: Reusable UI components

### Data Models
- `AppSettings`: App configuration data
- `ColorWithMeaning`: Color definitions with meanings
- `DayData`: Calendar day information

## 🔧 Configuration

### Minimum Requirements
- Min SDK: 35 (Android 15)
- Target SDK: 36
- Compile SDK: 36

### Dependencies
Key dependencies include:
- Jetpack Compose BOM
- Navigation Compose
- ViewModel Compose
- Kotlinx Serialization
- Material Icons Extended
- Kotlinx Coroutines

## 🧪 Testing

The project includes:
- Unit tests for core functionality
- Data model validation tests
- Color and date handling tests

Run tests with: `./gradlew test`

## 📦 Data Storage

The app uses local storage with:
- SharedPreferences for persistence
- JSON serialization for complex data
- Export/import functionality for data backup
- Graceful error handling for storage operations

## 🎨 Design System

- **Material 3** design language
- **Custom color schemes** for light/dark themes
- **Consistent typography** and spacing
- **Accessibility support** with semantic descriptions
- **Responsive layout** for different screen sizes

## 🚀 Future Enhancements

Potential improvements include:
- Cloud sync functionality
- Widget support
- Advanced analytics
- Reminder notifications
- Color themes and patterns
- Data visualization charts

## 📄 License

This project is developed for educational and personal use.

## 🤝 Contributing

This is a personal project, but suggestions and feedback are welcome!

---

Built with ❤️ using modern Android development practices.
