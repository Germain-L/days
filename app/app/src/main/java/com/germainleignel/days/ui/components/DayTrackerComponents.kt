package com.germainleignel.days.ui.components

import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyRow
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Check
import androidx.compose.material.icons.filled.Warning
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import com.germainleignel.days.data.ColorWithMeaning
import com.germainleignel.days.viewmodel.DayTrackerError
import kotlin.math.pow

// Rounded button component with 8dp radius, min-height 48dp
@Composable
fun DayTrackerButton(
    text: String,
    onClick: () -> Unit,
    modifier: Modifier = Modifier,
    enabled: Boolean = true
) {
    Button(
        onClick = onClick,
        modifier = modifier
            .fillMaxWidth()
            .heightIn(min = 48.dp),
        enabled = enabled,
        shape = RoundedCornerShape(8.dp),
        colors = ButtonDefaults.buttonColors(
            containerColor = MaterialTheme.colorScheme.primary,
            contentColor = MaterialTheme.colorScheme.onPrimary
        )
    ) {
        Text(
            text = text,
            style = MaterialTheme.typography.bodyLarge
        )
    }
}

// Toggle switch component with rounded track and circular thumb
@Composable
fun DayTrackerSwitch(
    checked: Boolean,
    onCheckedChange: (Boolean) -> Unit,
    modifier: Modifier = Modifier,
    enabled: Boolean = true
) {
    Switch(
        checked = checked,
        onCheckedChange = onCheckedChange,
        modifier = modifier,
        enabled = enabled,
        colors = SwitchDefaults.colors(
            checkedThumbColor = MaterialTheme.colorScheme.primary,
            checkedTrackColor = MaterialTheme.colorScheme.primary.copy(alpha = 0.5f),
            uncheckedThumbColor = MaterialTheme.colorScheme.onSurface,
            uncheckedTrackColor = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.3f)
        )
    )
}

// Card surface component with elevation and subtle shadow
@Composable
fun DayTrackerCard(
    modifier: Modifier = Modifier,
    content: @Composable () -> Unit
) {
    Card(
        modifier = modifier,
        shape = RoundedCornerShape(8.dp),
        elevation = CardDefaults.cardElevation(defaultElevation = 2.dp),
        colors = CardDefaults.cardColors(
            containerColor = MaterialTheme.colorScheme.surface
        )
    ) {
        content()
    }
}

// Divider component with 1dp line, full-width
@Composable
fun DayTrackerDivider(
    modifier: Modifier = Modifier
) {
    HorizontalDivider(
        modifier = modifier,
        thickness = 1.dp,
        color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.6f)
    )
}

// Color picker swatch component
@Composable
fun ColorSwatch(
    color: Color,
    isSelected: Boolean,
    onClick: () -> Unit,
    modifier: Modifier = Modifier
) {
    Box(
        modifier = modifier
            .size(48.dp)
            .clip(CircleShape)
            .background(color)
            .clickable { onClick() },
        contentAlignment = Alignment.Center
    ) {
        if (isSelected) {
            Icon(
                imageVector = Icons.Default.Check,
                contentDescription = "Selected",
                tint = if (color.luminance() > 0.5f) Color.Black else Color.White,
                modifier = Modifier.size(24.dp)
            )
        }
    }
}

// Extension function to calculate color luminance for accessibility
private fun Color.luminance(): Float {
    val r = if (red <= 0.03928f) red / 12.92f else (red + 0.055f).toDouble().pow(2.4).toFloat() / 1.055f.toDouble().pow(2.4).toFloat()
    val g = if (green <= 0.03928f) green / 12.92f else (green + 0.055f).toDouble().pow(2.4).toFloat() / 1.055f.toDouble().pow(2.4).toFloat()
    val b = if (blue <= 0.03928f) blue / 12.92f else (blue + 0.055f).toDouble().pow(2.4).toFloat() / 1.055f.toDouble().pow(2.4).toFloat()
    return 0.2126f * r + 0.7152f * g + 0.0722f * b
}

// Color meaning item for editing color meanings in settings
@Composable
fun ColorMeaningItem(
    color: Color,
    meaning: String,
    onMeaningChange: (String) -> Unit,
    modifier: Modifier = Modifier
) {
    var isEditing by remember { mutableStateOf(false) }
    var tempMeaning by remember(meaning) { mutableStateOf(meaning) }

    Row(
        modifier = modifier
            .fillMaxWidth()
            .padding(vertical = 8.dp),
        verticalAlignment = Alignment.CenterVertically,
        horizontalArrangement = Arrangement.spacedBy(12.dp)
    ) {
        // Color swatch
        ColorSwatch(
            color = color,
            isSelected = false,
            onClick = { isEditing = !isEditing },
            modifier = Modifier.size(40.dp)
        )

        // Meaning text or text field
        if (isEditing) {
            OutlinedTextField(
                value = tempMeaning,
                onValueChange = { tempMeaning = it },
                modifier = Modifier.weight(1f),
                singleLine = true,
                shape = RoundedCornerShape(8.dp),
                colors = OutlinedTextFieldDefaults.colors(
                    focusedBorderColor = MaterialTheme.colorScheme.primary,
                    unfocusedBorderColor = MaterialTheme.colorScheme.outline
                )
            )

            // Save button
            TextButton(
                onClick = {
                    onMeaningChange(tempMeaning)
                    isEditing = false
                }
            ) {
                Text("Save")
            }
        } else {
            // Display meaning
            Text(
                text = meaning,
                style = MaterialTheme.typography.bodyLarge,
                color = MaterialTheme.colorScheme.onSurface,
                modifier = Modifier.weight(1f)
            )

            // Edit button
            TextButton(onClick = { isEditing = true }) {
                Text("Edit")
            }
        }
    }
}

// Color management item for editing color meanings and deleting colors in settings
@Composable
fun ColorManagementItem(
    colorWithMeaning: ColorWithMeaning,
    onMeaningChange: (String) -> Unit,
    onDelete: () -> Unit,
    canDelete: Boolean,
    modifier: Modifier = Modifier
) {
    var isEditing by remember { mutableStateOf(false) }
    var tempMeaning by remember(colorWithMeaning.meaning) { mutableStateOf(colorWithMeaning.meaning) }

    Row(
        modifier = modifier
            .fillMaxWidth()
            .padding(vertical = 8.dp),
        verticalAlignment = Alignment.CenterVertically,
        horizontalArrangement = Arrangement.spacedBy(12.dp)
    ) {
        // Color swatch
        ColorSwatch(
            color = colorWithMeaning.color,
            isSelected = false,
            onClick = { isEditing = !isEditing },
            modifier = Modifier.size(40.dp)
        )

        // Meaning text or text field
        if (isEditing) {
            OutlinedTextField(
                value = tempMeaning,
                onValueChange = { tempMeaning = it },
                modifier = Modifier.weight(1f),
                singleLine = true,
                shape = RoundedCornerShape(8.dp),
                colors = OutlinedTextFieldDefaults.colors(
                    focusedBorderColor = MaterialTheme.colorScheme.primary,
                    unfocusedBorderColor = MaterialTheme.colorScheme.outline
                )
            )

            // Save button
            TextButton(
                onClick = {
                    onMeaningChange(tempMeaning)
                    isEditing = false
                }
            ) {
                Text("Save")
            }
        } else {
            // Display meaning
            Text(
                text = colorWithMeaning.meaning,
                style = MaterialTheme.typography.bodyLarge,
                color = MaterialTheme.colorScheme.onSurface,
                modifier = Modifier.weight(1f)
            )

            // Edit button
            TextButton(onClick = { isEditing = true }) {
                Text("Edit")
            }

            // Delete button (only if more than one color exists)
            if (canDelete) {
                TextButton(
                    onClick = onDelete,
                    colors = ButtonDefaults.textButtonColors(
                        contentColor = MaterialTheme.colorScheme.error
                    )
                ) {
                    Text("Delete")
                }
            }
        }
    }
}

// Add color dialog for creating new colors in settings
@Composable
fun AddColorDialog(
    onDismiss: () -> Unit,
    onAddColor: (Color, String) -> Unit
) {
    var selectedColor by remember { mutableStateOf(Color(0xFF9C27B0)) } // Default to purple
    var meaning by remember { mutableStateOf("") }

    // Predefined color options for selection
    val colorOptions = listOf(
        Color(0xFFE53E3E), // Red
        Color(0xFFE91E63), // Pink
        Color(0xFF9C27B0), // Purple
        Color(0xFF673AB7), // Deep Purple
        Color(0xFF3F51B5), // Indigo
        Color(0xFF2196F3), // Blue
        Color(0xFF00BCD4), // Cyan
        Color(0xFF009688), // Teal
        Color(0xFF4CAF50), // Green
        Color(0xFF8BC34A), // Light Green
        Color(0xFFCDDC39), // Lime
        Color(0xFFFFC107), // Amber
        Color(0xFFFF9800), // Orange
        Color(0xFFFF5722), // Deep Orange
        Color(0xFF795548), // Brown
        Color(0xFF607D8B), // Blue Grey
    )

    AlertDialog(
        onDismissRequest = onDismiss,
        title = {
            Text("Add New Color")
        },
        text = {
            Column(
                verticalArrangement = Arrangement.spacedBy(16.dp)
            ) {
                Text(
                    text = "Choose a color and give it a meaning",
                    style = MaterialTheme.typography.bodyMedium,
                    color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.7f)
                )

                // Color selection grid
                Text(
                    text = "Select Color:",
                    style = MaterialTheme.typography.bodyLarge,
                    color = MaterialTheme.colorScheme.onSurface
                )

                LazyRow(
                    horizontalArrangement = Arrangement.spacedBy(8.dp),
                    modifier = Modifier.padding(vertical = 8.dp)
                ) {
                    items(colorOptions.chunked(4)) { colorRow ->
                        Column(
                            verticalArrangement = Arrangement.spacedBy(8.dp)
                        ) {
                            colorRow.forEach { color ->
                                ColorSwatch(
                                    color = color,
                                    isSelected = color == selectedColor,
                                    onClick = { selectedColor = color },
                                    modifier = Modifier.size(40.dp)
                )
            }
        }
    }
}

// Error state component for displaying errors to users
@Composable
fun ErrorState(
    error: DayTrackerError,
    onRetry: () -> Unit,
    onDismiss: () -> Unit,
    modifier: Modifier = Modifier
) {
    Card(
        modifier = modifier.fillMaxWidth(),
        colors = CardDefaults.cardColors(
            containerColor = MaterialTheme.colorScheme.errorContainer
        ),
        shape = RoundedCornerShape(8.dp)
    ) {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(16.dp),
            horizontalArrangement = Arrangement.spacedBy(12.dp),
            verticalAlignment = Alignment.CenterVertically
        ) {
            Icon(
                imageVector = Icons.Default.Warning,
                contentDescription = "Error",
                tint = MaterialTheme.colorScheme.error,
                modifier = Modifier.size(24.dp)
            )

            Column(
                modifier = Modifier.weight(1f)
            ) {
                Text(
                    text = when (error) {
                        is DayTrackerError.StorageError -> "Storage Error"
                        is DayTrackerError.NetworkError -> "Network Error"
                        is DayTrackerError.ValidationError -> "Validation Error"
                        is DayTrackerError.UnknownError -> "Unknown Error"
                    },
                    style = MaterialTheme.typography.titleMedium,
                    color = MaterialTheme.colorScheme.error
                )

                Text(
                    text = when (error) {
                        is DayTrackerError.StorageError -> "Failed to save or load data. Please try again."
                        is DayTrackerError.NetworkError -> "Network connection failed. Check your internet connection."
                        is DayTrackerError.ValidationError -> error.message
                        is DayTrackerError.UnknownError -> "An unexpected error occurred. Please try again."
                    },
                    style = MaterialTheme.typography.bodyMedium,
                    color = MaterialTheme.colorScheme.onErrorContainer
                )
            }

            Row(
                horizontalArrangement = Arrangement.spacedBy(8.dp)
            ) {
                TextButton(onClick = onRetry) {
                    Text("Retry")
                }
                TextButton(onClick = onDismiss) {
                    Text("Dismiss")
                }
            }
        }
    }
}

// Loading state component
@Composable
fun LoadingState(
    message: String = "Loading...",
    modifier: Modifier = Modifier
) {
    Box(
        modifier = modifier.fillMaxWidth(),
        contentAlignment = Alignment.Center
    ) {
        Row(
            horizontalArrangement = Arrangement.spacedBy(12.dp),
            verticalAlignment = Alignment.CenterVertically,
            modifier = Modifier.padding(16.dp)
        ) {
            CircularProgressIndicator(
                modifier = Modifier.size(24.dp),
                strokeWidth = 2.dp
            )
            Text(
                text = message,
                style = MaterialTheme.typography.bodyMedium,
                color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.7f)
            )
        }
    }
}
                // Meaning input
                Text(
                    text = "Meaning:",
                    style = MaterialTheme.typography.bodyLarge,
                    color = MaterialTheme.colorScheme.onSurface
                )

                OutlinedTextField(
                    value = meaning,
                    onValueChange = { meaning = it },
                    placeholder = { Text("e.g., Productive, Relaxed, etc.") },
                    singleLine = true,
                    modifier = Modifier.fillMaxWidth(),
                    shape = RoundedCornerShape(8.dp)
                )
            }
        },
        confirmButton = {
            DayTrackerButton(
                text = "Add",
                onClick = {
                    if (meaning.isNotBlank()) {
                        onAddColor(selectedColor, meaning.trim())
                    }
                },
                enabled = meaning.isNotBlank(),
                modifier = Modifier.width(80.dp)
            )
        },
        dismissButton = {
            TextButton(onClick = onDismiss) {
                Text("Cancel")
            }
        }
    )
}
