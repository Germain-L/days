package com.germainleignel.days.ui.components

import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Add
import androidx.compose.material.icons.filled.Check
import androidx.compose.material.icons.filled.Delete
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import androidx.compose.ui.window.Dialog
import com.germainleignel.days.data.ColorWithMeaning
import com.germainleignel.days.data.getDefaultColors
import com.germainleignel.days.data.model.Calendar

@Composable
fun CreateCalendarDialog(
    onDismiss: () -> Unit,
    onCreate: (String, List<ColorWithMeaning>) -> Unit
) {
    var name by remember { mutableStateOf("") }
    var colors by remember { mutableStateOf(getDefaultColors()) }
    var showAddColorDialog by remember { mutableStateOf(false) }

    Dialog(onDismissRequest = onDismiss) {
        Card(
            modifier = Modifier
                .fillMaxWidth()
                .fillMaxHeight(0.8f),
            shape = RoundedCornerShape(16.dp)
        ) {
            Column(
                modifier = Modifier.padding(24.dp)
            ) {
                // Header
                Text(
                    text = "Create New Calendar",
                    style = MaterialTheme.typography.headlineSmall,
                    color = MaterialTheme.colorScheme.onSurface
                )

                Spacer(modifier = Modifier.height(24.dp))

                // Name input
                OutlinedTextField(
                    value = name,
                    onValueChange = { name = it },
                    label = { Text("Calendar Name") },
                    placeholder = { Text("e.g., Mental Health, Food Quality") },
                    modifier = Modifier.fillMaxWidth(),
                    singleLine = true
                )

                Spacer(modifier = Modifier.height(24.dp))

                // Colors section
                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.SpaceBetween,
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Text(
                        text = "Colors & Meanings",
                        style = MaterialTheme.typography.titleMedium,
                        color = MaterialTheme.colorScheme.onSurface
                    )
                    
                    IconButton(
                        onClick = { showAddColorDialog = true }
                    ) {
                        Icon(
                            imageVector = Icons.Default.Add,
                            contentDescription = "Add color",
                            tint = MaterialTheme.colorScheme.primary
                        )
                    }
                }

                Spacer(modifier = Modifier.height(12.dp))

                // Color list
                LazyColumn(
                    modifier = Modifier.weight(1f),
                    verticalArrangement = Arrangement.spacedBy(8.dp)
                ) {
                    items(colors) { colorWithMeaning ->
                        EditableColorItem(
                            colorWithMeaning = colorWithMeaning,
                            onMeaningChange = { newMeaning ->
                                colors = colors.map { 
                                    if (it.color == colorWithMeaning.color) {
                                        it.copy(meaning = newMeaning)
                                    } else it
                                }
                            },
                            onDelete = {
                                if (colors.size > 1) {
                                    colors = colors.filter { it.color != colorWithMeaning.color }
                                }
                            },
                            canDelete = colors.size > 1
                        )
                    }
                }

                Spacer(modifier = Modifier.height(24.dp))

                // Action buttons
                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.spacedBy(12.dp)
                ) {
                    TextButton(
                        onClick = onDismiss,
                        modifier = Modifier.weight(1f)
                    ) {
                        Text("Cancel")
                    }
                    
                    Button(
                        onClick = {
                            if (name.isNotBlank() && colors.isNotEmpty()) {
                                onCreate(name.trim(), colors)
                            }
                        },
                        enabled = name.isNotBlank() && colors.isNotEmpty(),
                        modifier = Modifier.weight(1f)
                    ) {
                        Text("Create")
                    }
                }
            }
        }
    }

    // Add color dialog
    if (showAddColorDialog) {
        AddColorDialog(
            onDismiss = { showAddColorDialog = false },
            onAddColor = { color, meaning ->
                colors = colors + ColorWithMeaning(color, meaning)
                showAddColorDialog = false
            }
        )
    }
}

@Composable
fun EditCalendarDialog(
    calendar: Calendar,
    onDismiss: () -> Unit,
    onSave: (Calendar) -> Unit
) {
    var name by remember { mutableStateOf(calendar.name) }
    var colors by remember { mutableStateOf(calendar.colorScheme) }
    var showAddColorDialog by remember { mutableStateOf(false) }

    Dialog(onDismissRequest = onDismiss) {
        Card(
            modifier = Modifier
                .fillMaxWidth()
                .fillMaxHeight(0.8f),
            shape = RoundedCornerShape(16.dp)
        ) {
            Column(
                modifier = Modifier.padding(24.dp)
            ) {
                // Header
                Text(
                    text = "Edit Calendar",
                    style = MaterialTheme.typography.headlineSmall,
                    color = MaterialTheme.colorScheme.onSurface
                )

                Spacer(modifier = Modifier.height(24.dp))

                // Name input
                OutlinedTextField(
                    value = name,
                    onValueChange = { name = it },
                    label = { Text("Calendar Name") },
                    modifier = Modifier.fillMaxWidth(),
                    singleLine = true
                )

                Spacer(modifier = Modifier.height(24.dp))

                // Colors section
                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.SpaceBetween,
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Text(
                        text = "Colors & Meanings",
                        style = MaterialTheme.typography.titleMedium,
                        color = MaterialTheme.colorScheme.onSurface
                    )
                    
                    IconButton(
                        onClick = { showAddColorDialog = true }
                    ) {
                        Icon(
                            imageVector = Icons.Default.Add,
                            contentDescription = "Add color",
                            tint = MaterialTheme.colorScheme.primary
                        )
                    }
                }

                Spacer(modifier = Modifier.height(12.dp))

                // Color list
                LazyColumn(
                    modifier = Modifier.weight(1f),
                    verticalArrangement = Arrangement.spacedBy(8.dp)
                ) {
                    items(colors) { colorWithMeaning ->
                        EditableColorItem(
                            colorWithMeaning = colorWithMeaning,
                            onMeaningChange = { newMeaning ->
                                colors = colors.map { 
                                    if (it.color == colorWithMeaning.color) {
                                        it.copy(meaning = newMeaning)
                                    } else it
                                }
                            },
                            onDelete = {
                                if (colors.size > 1) {
                                    colors = colors.filter { it.color != colorWithMeaning.color }
                                }
                            },
                            canDelete = colors.size > 1
                        )
                    }
                }

                Spacer(modifier = Modifier.height(24.dp))

                // Action buttons
                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.spacedBy(12.dp)
                ) {
                    TextButton(
                        onClick = onDismiss,
                        modifier = Modifier.weight(1f)
                    ) {
                        Text("Cancel")
                    }
                    
                    Button(
                        onClick = {
                            if (name.isNotBlank() && colors.isNotEmpty()) {
                                onSave(
                                    calendar.copy(
                                        name = name.trim(),
                                        colorScheme = colors
                                    )
                                )
                            }
                        },
                        enabled = name.isNotBlank() && colors.isNotEmpty(),
                        modifier = Modifier.weight(1f)
                    ) {
                        Text("Save")
                    }
                }
            }
        }
    }

    // Add color dialog
    if (showAddColorDialog) {
        AddColorDialog(
            onDismiss = { showAddColorDialog = false },
            onAddColor = { color, meaning ->
                colors = colors + ColorWithMeaning(color, meaning)
                showAddColorDialog = false
            }
        )
    }
}

@Composable
fun EditableColorItem(
    colorWithMeaning: ColorWithMeaning,
    onMeaningChange: (String) -> Unit,
    onDelete: () -> Unit,
    canDelete: Boolean
) {
    var isEditing by remember { mutableStateOf(false) }
    var tempMeaning by remember { mutableStateOf(colorWithMeaning.meaning) }

    Card(
        modifier = Modifier.fillMaxWidth(),
        colors = CardDefaults.cardColors(
            containerColor = MaterialTheme.colorScheme.surface
        )
    ) {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(16.dp),
            verticalAlignment = Alignment.CenterVertically,
            horizontalArrangement = Arrangement.spacedBy(12.dp)
        ) {
            // Color swatch
            Box(
                modifier = Modifier
                    .size(32.dp)
                    .clip(CircleShape)
                    .background(colorWithMeaning.color)
            )
            
            if (isEditing) {
                // Edit mode
                OutlinedTextField(
                    value = tempMeaning,
                    onValueChange = { tempMeaning = it },
                    modifier = Modifier.weight(1f),
                    singleLine = true
                )
                
                IconButton(
                    onClick = {
                        onMeaningChange(tempMeaning)
                        isEditing = false
                    }
                ) {
                    Icon(
                        imageVector = Icons.Default.Check,
                        contentDescription = "Save",
                        tint = MaterialTheme.colorScheme.primary
                    )
                }
            } else {
                // Display mode
                Text(
                    text = colorWithMeaning.meaning,
                    style = MaterialTheme.typography.bodyLarge,
                    color = MaterialTheme.colorScheme.onSurface,
                    modifier = Modifier
                        .weight(1f)
                        .clickable { 
                            tempMeaning = colorWithMeaning.meaning
                            isEditing = true 
                        }
                )
                
                if (canDelete) {
                    IconButton(onClick = onDelete) {
                        Icon(
                            imageVector = Icons.Default.Delete,
                            contentDescription = "Delete color",
                            tint = MaterialTheme.colorScheme.error
                        )
                    }
                }
            }
        }
    }
}
