package com.germainleignel.days.navigation

import androidx.compose.foundation.isSystemInDarkTheme
import androidx.compose.runtime.*
import androidx.compose.ui.platform.LocalContext
import androidx.lifecycle.viewmodel.compose.viewModel
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.rememberNavController
import com.germainleignel.days.api.auth.UserSessionManager
import com.germainleignel.days.storage.RepositoryFactory
import com.germainleignel.days.storage.StorageType
import com.germainleignel.days.ui.screens.AuthScreen
import com.germainleignel.days.ui.screens.CalendarScreen
import com.germainleignel.days.ui.screens.SettingsScreen
import com.germainleignel.days.ui.theme.DaysTheme
import com.germainleignel.days.viewmodel.DayTrackerViewModel

@Composable
fun DayTrackerApp() {
    val context = LocalContext.current
    val navController = rememberNavController()
    val viewModel: DayTrackerViewModel = viewModel()
    val settings by viewModel.settings.collectAsState()
    val isSystemInDarkTheme = isSystemInDarkTheme()

    // Check if we're using API mode and authentication status
    val currentStorageType = RepositoryFactory.getCurrentStorageType(context)
    val sessionManager = if (currentStorageType == StorageType.REMOTE) {
        RepositoryFactory.getSessionManager(context)
    } else null
    
    val authState by (sessionManager?.authState?.collectAsState() 
        ?: remember { mutableStateOf(UserSessionManager.AuthState.Authenticated) })

    // Determine if we should use dark theme
    val useDarkTheme = if (settings.followSystemTheme) {
        isSystemInDarkTheme
    } else {
        settings.isDarkMode
    }

    // Determine start destination based on auth state and storage type
    val startDestination = when {
        currentStorageType == StorageType.LOCAL -> "calendar"
        authState is UserSessionManager.AuthState.Authenticated -> "calendar"
        else -> "auth"
    }

    DaysTheme(darkTheme = useDarkTheme) {
        NavHost(
            navController = navController,
            startDestination = startDestination
        ) {
            composable("auth") {
                AuthScreen(
                    onAuthSuccess = {
                        navController.navigate("calendar") {
                            popUpTo("auth") { inclusive = true }
                        }
                    }
                )
            }

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
