package com.germainleignel.days.viewmodel

import android.content.Context
import androidx.lifecycle.ViewModel
import androidx.lifecycle.ViewModelProvider
import androidx.lifecycle.viewModelScope
import com.germainleignel.days.api.auth.UserSessionManager
import com.germainleignel.days.api.repository.ApiDataRepository
import com.germainleignel.days.storage.LocalDataRepository
import kotlinx.coroutines.launch

class AuthViewModel(
    private val sessionManager: UserSessionManager,
    context: Context
) : ViewModel() {
    
    private val apiRepository = ApiDataRepository(
        sessionManager,
        LocalDataRepository(context.applicationContext)
    )
    
    fun login(email: String, password: String) {
        viewModelScope.launch {
            apiRepository.login(email, password)
        }
    }
    
    fun register(email: String, password: String) {
        viewModelScope.launch {
            val result = apiRepository.createUser(email, password)
            if (result.isSuccess) {
                // After successful registration, automatically log in
                apiRepository.login(email, password)
            }
        }
    }
    
    fun logout() {
        apiRepository.logout()
    }
    
    class Factory(
        private val sessionManager: UserSessionManager,
        private val context: Context
    ) : ViewModelProvider.Factory {
        @Suppress("UNCHECKED_CAST")
        override fun <T : ViewModel> create(modelClass: Class<T>): T {
            if (modelClass.isAssignableFrom(AuthViewModel::class.java)) {
                return AuthViewModel(sessionManager, context) as T
            }
            throw IllegalArgumentException("Unknown ViewModel class")
        }
    }
}
