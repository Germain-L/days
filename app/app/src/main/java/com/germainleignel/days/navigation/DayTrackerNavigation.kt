package com.germainleignel.days.navigation

import androidx.compose.foundation.isSystemInDarkTheme
import androidx.compose.runtime.*
import androidx.lifecycle.viewmodel.compose.viewModel
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.rememberNavController
import com.germainleignel.days.ui.screens.CalendarScreen
import com.germainleignel.days.ui.screens.SettingsScreen
import com.germainleignel.days.ui.theme.DaysTheme
import com.germainleignel.days.viewmodel.DayTrackerViewModel

@Composable
fun DayTrackerApp() {
    val navController = rememberNavController()
    val viewModel: DayTrackerViewModel = viewModel()
    val settings by viewModel.settings.collectAsState()
    val isSystemInDarkTheme = isSystemInDarkTheme()

    // Determine if we should use dark theme
    val useDarkTheme = if (settings.followSystemTheme) {
        isSystemInDarkTheme
    } else {
        settings.isDarkMode
    }

    DaysTheme(darkTheme = useDarkTheme) {
        NavHost(
            navController = navController,
            startDestination = "calendar"
        ) {
            composable("calendar") {
                CalendarScreen(
                    onNavigateToSettings = {
                        navController.navigate("settings")
                    },
                    viewModel = viewModel
                )
            }

            composable("settings") {
                SettingsScreen(
                    onNavigateBack = {
                        navController.popBackStack()
                    },
                    viewModel = viewModel
                )
            }
        }
    }
}
