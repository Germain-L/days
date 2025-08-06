package com.germainleignel.days.data.model

import com.germainleignel.days.data.ColorWithMeaning
import java.util.UUID

data class Calendar(
    val id: String = UUID.randomUUID().toString(),
    val name: String,
    val colorScheme: List<ColorWithMeaning>,
    val isSelected: Boolean = false,
    val createdAt: Long = System.currentTimeMillis()
) {
    companion object {
        fun createDefault(name: String, colors: List<ColorWithMeaning>): Calendar {
            return Calendar(
                name = name,
                colorScheme = colors,
                isSelected = true
            )
        }
    }
}
