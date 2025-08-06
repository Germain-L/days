package com.germainleignel.days.ui.screens

import androidx.compose.foundation.background
import androidx.compose.foundation.isSystemInDarkTheme
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyRow
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.ArrowBack
import androidx.compose.material.icons.filled.FileDownload
import androidx.compose.material.icons.filled.FileUpload
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.LocalClipboardManager
import androidx.compose.ui.text.AnnotatedString
import androidx.compose.ui.unit.dp
import androidx.lifecycle.viewmodel.compose.viewModel
import com.germainleignel.days.ui.components.*
import com.germainleignel.days.viewmodel.DayTrackerViewModel
import com.germainleignel.days.data.ColorWithMeaning

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun SettingsScreen(
    onNavigateBack: () -> Unit,
    viewModel: DayTrackerViewModel = viewModel()
) {
    val settings by viewModel.settings.collectAsState()
    val isSystemInDarkTheme = isSystemInDarkTheme()
    var showResetDialog by remember { mutableStateOf(false) }
    var showColorBottomSheet by remember { mutableStateOf(false) }
    var showAddColorDialog by remember { mutableStateOf(false) }
    var showExportDialog by remember { mutableStateOf(false) }
    var showImportDialog by remember { mutableStateOf(false) }
    var exportedData by remember { mutableStateOf<String?>(null) }
    var importText by remember { mutableStateOf("") }

    val clipboardManager = LocalClipboardManager.current

    val bottomSheetState = rememberModalBottomSheetState()

    Column(
        modifier = Modifier
            .fillMaxSize()
            .background(MaterialTheme.colorScheme.background)
    ) {
        // Top app bar
        TopAppBar(
            title = {
                Text(
                    text = "Settings",
                    style = MaterialTheme.typography.titleLarge
                )
            },
            navigationIcon = {
                IconButton(onClick = onNavigateBack) {
                    Icon(
                        imageVector = Icons.Default.ArrowBack,
                        contentDescription = "Back"
                    )
                }
            },
            colors = TopAppBarDefaults.topAppBarColors(
                containerColor = MaterialTheme.colorScheme.background
            )
        )

        // Settings content
        Column(
            modifier = Modifier
                .fillMaxSize()
                .verticalScroll(rememberScrollState())
                .padding(24.dp),
            verticalArrangement = Arrangement.spacedBy(32.dp)
        ) {
            // Color selection section
            DayTrackerCard {
                Column(
                    modifier = Modifier.padding(20.dp)
                ) {
                    Text(
                        text = "Default Color",
                        style = MaterialTheme.typography.titleLarge,
                        color = MaterialTheme.colorScheme.onSurface,
                        modifier = Modifier.padding(bottom = 20.dp)
                    )

                    Text(
                        text = "Choose the color that will be applied when you tap a day",
                        style = MaterialTheme.typography.bodyLarge,
                        color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.7f),
                        modifier = Modifier.padding(bottom = 20.dp)
                    )

                    // Current selected color display with click to open bottom sheet
                    Row(
                        modifier = Modifier
                            .fillMaxWidth()
                            .padding(bottom = 20.dp),
                        verticalAlignment = Alignment.CenterVertically,
                        horizontalArrangement = Arrangement.SpaceBetween
                    ) {
                        Row(
                            verticalAlignment = Alignment.CenterVertically,
                            horizontalArrangement = Arrangement.spacedBy(16.dp)
                        ) {
                            ColorSwatch(
                                color = settings.selectedColor,
                                isSelected = true,
                                onClick = { showColorBottomSheet = true }
                            )
                            Column {
                                Text(
                                    text = "Current Color",
                                    style = MaterialTheme.typography.bodyLarge,
                                    color = MaterialTheme.colorScheme.onSurface
                                )
                                Text(
                                    text = viewModel.getColorMeaning(settings.selectedColor),
                                    style = MaterialTheme.typography.bodyMedium,
                                    color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.7f)
                                )
                            }
                        }

                        DayTrackerButton(
                            text = "Change Color",
                            onClick = { showColorBottomSheet = true },
                            modifier = Modifier.width(120.dp)
                        )
                    }
                }
            }

            // Color meanings section
            DayTrackerCard {
                Column(
                    modifier = Modifier.padding(20.dp)
                ) {
                    Row(
                        modifier = Modifier.fillMaxWidth(),
                        horizontalArrangement = Arrangement.SpaceBetween,
                        verticalAlignment = Alignment.CenterVertically
                    ) {
                        Text(
                            text = "Manage Colors",
                            style = MaterialTheme.typography.titleLarge,
                            color = MaterialTheme.colorScheme.onSurface
                        )

                        DayTrackerButton(
                            text = "Add Color",
                            onClick = { showAddColorDialog = true },
                            modifier = Modifier.width(100.dp)
                        )
                    }

                    Text(
                        text = "Add, edit, or remove colors and their meanings",
                        style = MaterialTheme.typography.bodyLarge,
                        color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.7f),
                        modifier = Modifier.padding(top = 8.dp, bottom = 20.dp)
                    )

                    // List of colors with their meanings and delete option
                    Column(
                        verticalArrangement = Arrangement.spacedBy(16.dp)
                    ) {
                        viewModel.getAvailableColorsWithMeanings().forEach { colorWithMeaning ->
                            ColorManagementItem(
                                colorWithMeaning = colorWithMeaning,
                                onMeaningChange = { newMeaning ->
                                    viewModel.updateColorMeaning(colorWithMeaning.color, newMeaning)
                                },
                                onDelete = {
                                    if (viewModel.getAvailableColors().size > 1) {
                                        viewModel.deleteColor(colorWithMeaning.color)
                                    }
                                },
                                canDelete = viewModel.getAvailableColors().size > 1
                            )
                        }
                    }
                }
            }

            // Theme section
            DayTrackerCard {
                Column(
                    modifier = Modifier.padding(20.dp)
                ) {
                    Text(
                        text = "Appearance",
                        style = MaterialTheme.typography.titleLarge,
                        color = MaterialTheme.colorScheme.onSurface,
                        modifier = Modifier.padding(bottom = 20.dp)
                    )

                    // Follow system theme option
                    Row(
                        modifier = Modifier
                            .fillMaxWidth()
                            .padding(bottom = 20.dp),
                        horizontalArrangement = Arrangement.SpaceBetween,
                        verticalAlignment = Alignment.CenterVertically
                    ) {
                        Column(modifier = Modifier.weight(1f)) {
                            Text(
                                text = "Follow System Theme",
                                style = MaterialTheme.typography.bodyLarge,
                                color = MaterialTheme.colorScheme.onSurface
                            )
                            Text(
                                text = "Automatically switch between light and dark themes",
                                style = MaterialTheme.typography.bodyLarge,
                                color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.7f)
                            )
                        }

                        DayTrackerSwitch(
                            checked = settings.followSystemTheme,
                            onCheckedChange = { follow ->
                                viewModel.setFollowSystemTheme(follow)
                                if (follow) {
                                    viewModel.updateDarkMode(isSystemInDarkTheme)
                                }
                            }
                        )
                    }

                    DayTrackerDivider()

                    Spacer(modifier = Modifier.height(20.dp))

                    // Manual dark mode toggle (only when not following system)
                    Row(
                        modifier = Modifier.fillMaxWidth(),
                        horizontalArrangement = Arrangement.SpaceBetween,
                        verticalAlignment = Alignment.CenterVertically
                    ) {
                        Column(modifier = Modifier.weight(1f)) {
                            Text(
                                text = "Dark Mode",
                                style = MaterialTheme.typography.bodyLarge,
                                color = if (settings.followSystemTheme) {
                                    MaterialTheme.colorScheme.onSurface.copy(alpha = 0.5f)
                                } else {
                                    MaterialTheme.colorScheme.onSurface
                                }
                            )
                            Text(
                                text = if (settings.followSystemTheme) {
                                    "Controlled by system settings"
                                } else {
                                    "Use dark theme colors"
                                },
                                style = MaterialTheme.typography.bodyLarge,
                                color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.7f)
                            )
                        }

                        DayTrackerSwitch(
                            checked = if (settings.followSystemTheme) isSystemInDarkTheme else settings.isDarkMode,
                            onCheckedChange = { isDark ->
                                if (!settings.followSystemTheme) {
                                    viewModel.updateDarkMode(isDark)
                                }
                            },
                            enabled = !settings.followSystemTheme
                        )
                    }
                }
            }

            // Data management section
            DayTrackerCard {
                Column(
                    modifier = Modifier.padding(20.dp)
                ) {
                    Text(
                        text = "Data Management",
                        style = MaterialTheme.typography.titleLarge,
                        color = MaterialTheme.colorScheme.onSurface,
                        modifier = Modifier.padding(bottom = 20.dp)
                    )

                    Text(
                        text = "Export your data for backup or import from a previous backup",
                        style = MaterialTheme.typography.bodyLarge,
                        color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.7f),
                        modifier = Modifier.padding(bottom = 20.dp)
                    )

                    // Export/Import buttons
                    Row(
                        modifier = Modifier.fillMaxWidth(),
                        horizontalArrangement = Arrangement.spacedBy(12.dp)
                    ) {
                        DayTrackerButton(
                            text = "Export Data",
                            onClick = { 
                                viewModel.exportData { data ->
                                    exportedData = data
                                    showExportDialog = true
                                }
                            },
                            modifier = Modifier.weight(1f)
                        )

                        DayTrackerButton(
                            text = "Import Data",
                            onClick = { showImportDialog = true },
                            modifier = Modifier.weight(1f)
                        )
                    }

                    Spacer(modifier = Modifier.height(20.dp))

                    DayTrackerDivider()

                    Spacer(modifier = Modifier.height(20.dp))

                    Text(
                        text = "Reset all colored days and settings to defaults",
                        style = MaterialTheme.typography.bodyLarge,
                        color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.7f),
                        modifier = Modifier.padding(bottom = 20.dp)
                    )

                    DayTrackerButton(
                        text = "Reset All Data",
                        onClick = { showResetDialog = true }
                    )
                }
            }

            // App info section
            DayTrackerCard {
                Column(
                    modifier = Modifier.padding(20.dp)
                ) {
                    Text(
                        text = "About",
                        style = MaterialTheme.typography.titleLarge,
                        color = MaterialTheme.colorScheme.onSurface,
                        modifier = Modifier.padding(bottom = 20.dp)
                    )

                    Text(
                        text = "Day Tracker",
                        style = MaterialTheme.typography.bodyLarge,
                        color = MaterialTheme.colorScheme.onSurface,
                        modifier = Modifier.padding(bottom = 4.dp)
                    )

                    Text(
                        text = "Version 1.0.0",
                        style = MaterialTheme.typography.bodyLarge,
                        color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.7f),
                        modifier = Modifier.padding(bottom = 8.dp)
                    )

                    Text(
                        text = "Track your days with colors. Tap any day to mark it with your selected color, or long-press for more color options.",
                        style = MaterialTheme.typography.bodyLarge,
                        color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.7f)
                    )
                }
            }
        }
    }

    // Reset confirmation dialog
    if (showResetDialog) {
        AlertDialog(
            onDismissRequest = { showResetDialog = false },
            title = {
                Text("Reset All Data")
            },
            text = {
                Text("This will remove all colored days and reset settings to defaults. This action cannot be undone.")
            },
            confirmButton = {
                TextButton(
                    onClick = {
                        viewModel.resetAllData()
                        showResetDialog = false
                    }
                ) {
                    Text(
                        "Reset",
                        color = MaterialTheme.colorScheme.error
                    )
                }
            },
            dismissButton = {
                TextButton(onClick = { showResetDialog = false }) {
                    Text("Cancel")
                }
            }
        )
    }

    // Color selection bottom sheet
    if (showColorBottomSheet) {
        ModalBottomSheet(
            onDismissRequest = { showColorBottomSheet = false },
            sheetState = bottomSheetState
        ) {
            Column(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(24.dp),
                verticalArrangement = Arrangement.spacedBy(20.dp)
            ) {
                Text(
                    text = "Select a Color",
                    style = MaterialTheme.typography.titleLarge,
                    color = MaterialTheme.colorScheme.onSurface
                )

                // Color options
                LazyRow(
                    horizontalArrangement = Arrangement.spacedBy(16.dp)
                ) {
                    items(viewModel.getAvailableColors()) { color ->
                        ColorSwatch(
                            color = color,
                            isSelected = color == settings.selectedColor,
                            onClick = {
                                viewModel.updateSelectedColor(color)
                                showColorBottomSheet = false
                            }
                        )
                    }
                }

                // Close button
                TextButton(
                    onClick = { showColorBottomSheet = false },
                    modifier = Modifier.align(Alignment.End)
                ) {
                    Text("Close")
                }
            }
        }
    }

    // Add color dialog
    if (showAddColorDialog) {
        AddColorDialog(
            onDismiss = { showAddColorDialog = false },
            onAddColor = { color, meaning ->
                viewModel.addNewColor(color, meaning)
                showAddColorDialog = false
            }
        )
    }

    // Export data dialog
    if (showExportDialog && exportedData != null) {
        AlertDialog(
            onDismissRequest = { 
                showExportDialog = false 
                exportedData = null
            },
            title = { Text("Export Data") },
            text = {
                Column {
                    Text("Your data has been exported. Copy the text below and save it securely:")
                    Spacer(modifier = Modifier.height(8.dp))
                    OutlinedTextField(
                        value = exportedData!!,
                        onValueChange = { },
                        readOnly = true,
                        modifier = Modifier
                            .fillMaxWidth()
                            .height(120.dp),
                        textStyle = MaterialTheme.typography.bodySmall
                    )
                }
            },
            confirmButton = {
                TextButton(
                    onClick = {
                        clipboardManager.setText(AnnotatedString(exportedData!!))
                        showExportDialog = false
                        exportedData = null
                    }
                ) {
                    Text("Copy to Clipboard")
                }
            },
            dismissButton = {
                TextButton(
                    onClick = { 
                        showExportDialog = false 
                        exportedData = null
                    }
                ) {
                    Text("Close")
                }
            }
        )
    }

    // Import data dialog
    if (showImportDialog) {
        AlertDialog(
            onDismissRequest = { 
                showImportDialog = false 
                importText = ""
            },
            title = { Text("Import Data") },
            text = {
                Column {
                    Text("Paste your exported data below. This will replace all current data.")
                    Spacer(modifier = Modifier.height(8.dp))
                    OutlinedTextField(
                        value = importText,
                        onValueChange = { importText = it },
                        placeholder = { Text("Paste exported data here...") },
                        modifier = Modifier
                            .fillMaxWidth()
                            .height(120.dp),
                        textStyle = MaterialTheme.typography.bodySmall
                    )
                }
            },
            confirmButton = {
                TextButton(
                    onClick = {
                        if (importText.isNotBlank()) {
                            viewModel.importData(importText) { success ->
                                // Handle import result
                                showImportDialog = false
                                importText = ""
                            }
                        }
                    },
                    enabled = importText.isNotBlank()
                ) {
                    Text("Import")
                }
            },
            dismissButton = {
                TextButton(
                    onClick = { 
                        showImportDialog = false 
                        importText = ""
                    }
                ) {
                    Text("Cancel")
                }
            }
        )
    }
}
