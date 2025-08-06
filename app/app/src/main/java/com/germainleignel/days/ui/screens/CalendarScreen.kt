package com.germainleignel.days.ui.screens

import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.LazyRow
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.lazy.rememberLazyListState
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.KeyboardArrowLeft
import androidx.compose.material.icons.filled.KeyboardArrowRight
import androidx.compose.material.icons.filled.Settings
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.semantics.contentDescription
import androidx.compose.ui.semantics.semantics
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import androidx.lifecycle.viewmodel.compose.viewModel
import com.germainleignel.days.data.ColorPalette
import com.germainleignel.days.ui.components.ColorSwatch
import com.germainleignel.days.ui.components.CalendarSelector
import com.germainleignel.days.ui.components.CalendarManagementDialog
import com.germainleignel.days.viewmodel.DayTrackerViewModel
import java.time.LocalDate
import java.time.YearMonth
import java.time.format.DateTimeFormatter
import java.time.format.TextStyle
import java.util.*
import kotlin.math.pow

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun CalendarScreen(
    onNavigateToSettings: () -> Unit,
    viewModel: DayTrackerViewModel = viewModel(),
) {
    val currentMonth = remember { YearMonth.now() }
    var displayedMonth by remember { mutableStateOf(currentMonth) }
    val listState = rememberLazyListState(
        initialFirstVisibleItemIndex = 0 // Start at current month
    )

    var showColorBottomSheet by remember { mutableStateOf(false) }
    var selectedDate by remember { mutableStateOf<LocalDate?>(null) }
    val bottomSheetState = rememberModalBottomSheetState()
    var showCalendarManagement by remember { mutableStateOf(false) }

    val isLoading: Boolean by viewModel.isLoading.collectAsState()
    val calendarData by viewModel.calendarData.collectAsState()
    val error: com.germainleignel.days.viewmodel.DayTrackerError? by viewModel.error.collectAsState()

    Column(
        modifier = Modifier
            .fillMaxSize()
            .background(MaterialTheme.colorScheme.background)
    ) {
        // Top app bar with calendar selector and settings button
        TopAppBar(
            title = {
                CalendarSelector(
                    calendars = calendarData.calendars,
                    selectedCalendar = calendarData.getSelectedCalendar(),
                    onCalendarSelected = { calendar ->
                        viewModel.selectCalendar(calendar.id)
                    },
                    onCreateCalendar = {
                        showCalendarManagement = true
                    },
                    modifier = Modifier.weight(1f)
                )
            },
            actions = {
                IconButton(onClick = onNavigateToSettings) {
                    Icon(
                        imageVector = Icons.Default.Settings,
                        contentDescription = "Settings"
                    )
                }
            },
            colors = TopAppBarDefaults.topAppBarColors(
                containerColor = MaterialTheme.colorScheme.background
            )
        )

        // TODO: Add error and loading states back after fixing compilation issues

        // Infinite scroll calendar
        LazyColumn(
            state = listState,
            modifier = Modifier.fillMaxSize(),
            contentPadding = PaddingValues(24.dp),
            verticalArrangement = Arrangement.spacedBy(40.dp)
        ) {
            items(60) { monthOffset -> // Show 5 years worth of months (60 months)
                val month = currentMonth.plusMonths(monthOffset.toLong())
                MonthSection(
                    month = month,
                    viewModel = viewModel,
                    onDayClick = { date ->
                        viewModel.toggleDayColor(date)
                    },
                    onDayLongClick = { date ->
                        selectedDate = date
                        showColorBottomSheet = true
                    },
                    onPreviousMonth = {
                        displayedMonth = displayedMonth.minusMonths(1)
                    },
                    onNextMonth = {
                        displayedMonth = displayedMonth.plusMonths(1)
                    }
                )
            }
        }
    }

    // Color picker bottom sheet
    if (showColorBottomSheet && selectedDate != null) {
        ColorPickerBottomSheet(
            selectedDate = selectedDate!!,
            currentColor = viewModel.getDayColor(selectedDate!!),
            onColorSelected = { color ->
                viewModel.setDayColor(selectedDate!!, color)
                showColorBottomSheet = false
                selectedDate = null
            },
            onDismiss = {
                showColorBottomSheet = false
                selectedDate = null
            }
        )
    }

    // Calendar management dialog
    if (showCalendarManagement) {
        CalendarManagementDialog(
            calendars = calendarData.calendars,
            onDismiss = { showCalendarManagement = false },
            onCreateCalendar = { name, colors ->
                viewModel.createCalendar(name, colors)
            },
            onEditCalendar = { calendar ->
                viewModel.updateCalendar(calendar)
            },
            onDeleteCalendar = { calendarId ->
                viewModel.deleteCalendar(calendarId)
            }
        )
    }
}

@Composable
fun MonthSection(
    month: YearMonth,
    viewModel: DayTrackerViewModel,
    onDayClick: (LocalDate) -> Unit,
    onDayLongClick: (LocalDate) -> Unit,
    onPreviousMonth: () -> Unit,
    onNextMonth: () -> Unit,
) {
    Column {
        // Month header
        Row(
            modifier = Modifier.fillMaxWidth(),
            horizontalArrangement = Arrangement.SpaceBetween,
            verticalAlignment = Alignment.CenterVertically
        ) {
            IconButton(onClick = onPreviousMonth) {
                Icon(
                    imageVector = Icons.Default.KeyboardArrowLeft,
                    contentDescription = "Previous month"
                )
            }

            Text(
                text = month.format(DateTimeFormatter.ofPattern("MMMM yyyy")),
                style = MaterialTheme.typography.headlineLarge,
                color = MaterialTheme.colorScheme.onBackground
            )

            IconButton(onClick = onNextMonth) {
                Icon(
                    imageVector = Icons.Default.KeyboardArrowRight,
                    contentDescription = "Next month"
                )
            }
        }

        Spacer(modifier = Modifier.height(24.dp))

        // Days of week header
        Row(
            modifier = Modifier.fillMaxWidth(),
            horizontalArrangement = Arrangement.SpaceEvenly
        ) {
            listOf("Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun").forEach { dayName ->
                Text(
                    text = dayName,
                    style = MaterialTheme.typography.labelMedium,
                    color = MaterialTheme.colorScheme.onBackground.copy(alpha = 0.6f),
                    modifier = Modifier.weight(1f),
                    textAlign = TextAlign.Center
                )
            }
        }

        Spacer(modifier = Modifier.height(16.dp))

        // Calendar grid
        val firstDayOfMonth = month.atDay(1)
        val daysInMonth = month.lengthOfMonth()
        val firstDayOfWeek = firstDayOfMonth.dayOfWeek.value - 1 // Monday = 0

        // Create weeks
        val weeks = mutableListOf<List<LocalDate?>>()
        var currentWeek = mutableListOf<LocalDate?>()

        // Add empty days at the beginning
        repeat(firstDayOfWeek) {
            currentWeek.add(null)
        }

        // Add all days of the month
        for (day in 1..daysInMonth) {
            if (currentWeek.size == 7) {
                weeks.add(currentWeek.toList())
                currentWeek.clear()
            }
            currentWeek.add(month.atDay(day))
        }

        // Add empty days at the end
        while (currentWeek.size < 7) {
            currentWeek.add(null)
        }
        weeks.add(currentWeek.toList())

        // Render weeks
        Column(
            verticalArrangement = Arrangement.spacedBy(12.dp)
        ) {
            weeks.forEach { week ->
                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.SpaceEvenly
                ) {
                    week.forEach { date ->
                        DayTile(
                            date = date,
                            color = date?.let { viewModel.getDayColor(it) },
                            onClick = { date?.let { onDayClick(it) } },
                            onLongClick = { date?.let { onDayLongClick(it) } },
                            modifier = Modifier.weight(1f)
                        )
                    }
                }
            }
        }
    }
}

@Composable
fun DayTile(
    date: LocalDate?,
    color: Color?,
    onClick: () -> Unit,
    onLongClick: () -> Unit,
    modifier: Modifier = Modifier,
) {
    Box(
        modifier = modifier
            .aspectRatio(1f)
            .padding(6.dp)
            .border(
                width = 2.dp,
                color = if (color != null) MaterialTheme.colorScheme.primary else MaterialTheme.colorScheme.outline.copy(
                    alpha = 0.3f
                ),
                shape = RoundedCornerShape(12.dp)
            )
            .clip(RoundedCornerShape(12.dp))
            .background(
                color ?: MaterialTheme.colorScheme.surface
            )
            .clickable(enabled = date != null) { onClick() }
            .then(
                if (date != null) {
                    Modifier
                        .clickable(
                            onClickLabel = "Long press for color picker"
                        ) { onLongClick() }
                        .semantics {
                            contentDescription = if (color != null) {
                                "Day ${date.dayOfMonth}, colored"
                            } else {
                                "Day ${date.dayOfMonth}, not colored"
                            }
                        }
                } else Modifier
            ),
        contentAlignment = Alignment.Center
    ) {
        if (date != null) {
            Text(
                text = date.dayOfMonth.toString(),
                style = MaterialTheme.typography.titleMedium,
                color = if (color != null) {
                    if (color.luminance() > 0.5f) Color.Black else Color.White
                } else {
                    MaterialTheme.colorScheme.onSurface
                },
                textAlign = TextAlign.Center
            )
        }
    }
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun ColorPickerBottomSheet(
    selectedDate: LocalDate,
    currentColor: Color?,
    onColorSelected: (Color) -> Unit,
    onDismiss: () -> Unit,
) {
    val viewModel: DayTrackerViewModel = viewModel()

    ModalBottomSheet(
        onDismissRequest = onDismiss,
        sheetState = rememberModalBottomSheetState()
    ) {
        Column(
            modifier = Modifier
                .fillMaxWidth()
                .padding(32.dp),
            horizontalAlignment = Alignment.CenterHorizontally
        ) {
            Text(
                text = "Choose Color",
                style = MaterialTheme.typography.titleLarge,
                modifier = Modifier.padding(bottom = 12.dp)
            )

            Text(
                text = selectedDate.format(DateTimeFormatter.ofPattern("MMMM d, yyyy")),
                style = MaterialTheme.typography.bodyLarge,
                color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.7f),
                modifier = Modifier.padding(bottom = 32.dp)
            )

            // Color options with meanings
            Column(
                verticalArrangement = Arrangement.spacedBy(24.dp),
                modifier = Modifier.padding(bottom = 32.dp)
            ) {
                // Add "Remove Color" option at the top if there's currently a color
                if (currentColor != null) {
                    Row(
                        horizontalArrangement = Arrangement.Center,
                        modifier = Modifier.fillMaxWidth()
                    ) {
                        Column(
                            horizontalAlignment = Alignment.CenterHorizontally,
                            modifier = Modifier
                                .clickable { onColorSelected(Color.Transparent) }
                                .padding(8.dp)
                        ) {
                            Box(
                                modifier = Modifier
                                    .size(56.dp)
                                    .clip(RoundedCornerShape(28.dp))
                                    .background(MaterialTheme.colorScheme.errorContainer)
                                    .border(
                                        2.dp,
                                        MaterialTheme.colorScheme.error,
                                        RoundedCornerShape(28.dp)
                                    ),
                                contentAlignment = Alignment.Center
                            ) {
                                Text(
                                    text = "×",
                                    style = MaterialTheme.typography.headlineLarge,
                                    color = MaterialTheme.colorScheme.error
                                )
                            }

                            Spacer(modifier = Modifier.height(8.dp))

                            Text(
                                text = "Remove Color",
                                style = MaterialTheme.typography.bodySmall,
                                color = MaterialTheme.colorScheme.error,
                                textAlign = TextAlign.Center,
                                maxLines = 2
                            )
                        }
                    }

                    // Divider
                    HorizontalDivider(
                        modifier = Modifier.padding(vertical = 16.dp),
                        color = MaterialTheme.colorScheme.outline.copy(alpha = 0.3f)
                    )
                }

                // Color palette in grid
                ColorPalette.chunked(3).forEach { colorRow ->
                    Row(
                        horizontalArrangement = Arrangement.SpaceEvenly,
                        modifier = Modifier.fillMaxWidth()
                    ) {
                        colorRow.forEach { color ->
                            Column(
                                horizontalAlignment = Alignment.CenterHorizontally,
                                modifier = Modifier
                                    .weight(1f)
                                    .clickable { onColorSelected(color) }
                                    .padding(12.dp)
                            ) {
                                ColorSwatch(
                                    color = color,
                                    isSelected = color == currentColor,
                                    onClick = { onColorSelected(color) },
                                    modifier = Modifier.size(64.dp)
                                )

                                Spacer(modifier = Modifier.height(12.dp))

                                Text(
                                    text = viewModel.getColorMeaning(color),
                                    style = MaterialTheme.typography.bodySmall,
                                    color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.8f),
                                    textAlign = TextAlign.Center,
                                    maxLines = 2
                                )
                            }
                        }

                        // Fill empty spaces in incomplete rows
                        repeat(3 - colorRow.size) {
                            Spacer(modifier = Modifier.weight(1f))
                        }
                    }
                }
            }

            // Action buttons
            Row(
                horizontalArrangement = Arrangement.spacedBy(16.dp),
                modifier = Modifier.padding(bottom = 16.dp)
            ) {
                TextButton(onClick = onDismiss) {
                    Text("Cancel")
                }
            }
        }
    }
}

// Extension function for color luminance calculation
private fun Color.luminance(): Float {
    val r = if (red <= 0.03928f) red / 12.92f else (red + 0.055f).toDouble().pow(2.4)
        .toFloat() / 1.055f.toDouble().pow(2.4).toFloat()
    val g = if (green <= 0.03928f) green / 12.92f else (green + 0.055f).toDouble().pow(2.4)
        .toFloat() / 1.055f.toDouble().pow(2.4).toFloat()
    val b = if (blue <= 0.03928f) blue / 12.92f else (blue + 0.055f).toDouble().pow(2.4)
        .toFloat() / 1.055f.toDouble().pow(2.4).toFloat()
    return 0.2126f * r + 0.7152f * g + 0.0722f * b
}
