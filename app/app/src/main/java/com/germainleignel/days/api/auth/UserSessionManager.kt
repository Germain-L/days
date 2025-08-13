package com.germainleignel.days.api.auth

import android.content.Context
import androidx.security.crypto.EncryptedSharedPreferences
import androidx.security.crypto.MasterKeys
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow

/**
 * Manages user authentication state and JWT tokens
 */
class UserSessionManager(context: Context) {
    private val masterKeyAlias = MasterKeys.getOrCreate(MasterKeys.AES256_GCM_SPEC)
    
    private val sharedPreferences = EncryptedSharedPreferences.create(
        "user_session",
        masterKeyAlias,
        context,
        EncryptedSharedPreferences.PrefKeyEncryptionScheme.AES256_SIV,
        EncryptedSharedPreferences.PrefValueEncryptionScheme.AES256_GCM
    )

    private val _authState = MutableStateFlow<AuthState>(AuthState.Unauthenticated)
    val authState: StateFlow<AuthState> = _authState.asStateFlow()

    private val _currentUser = MutableStateFlow<User?>(null)
    val currentUser: StateFlow<User?> = _currentUser.asStateFlow()

    data class User(
        val id: String,
        val email: String,
        val createdAt: String
    )

    sealed class AuthState {
        object Loading : AuthState()
        object Authenticated : AuthState()
        object Unauthenticated : AuthState()
        data class Error(val message: String) : AuthState()
    }

    companion object {
        private const val KEY_JWT_TOKEN = "jwt_token"
        private const val KEY_USER_ID = "user_id"
        private const val KEY_USER_EMAIL = "user_email"
        private const val KEY_USER_CREATED_AT = "user_created_at"
    }

    init {
        // Check if user is already logged in
        if (hasValidSession()) {
            _authState.value = AuthState.Authenticated
            _currentUser.value = User(
                id = sharedPreferences.getString(KEY_USER_ID, "") ?: "",
                email = sharedPreferences.getString(KEY_USER_EMAIL, "") ?: "",
                createdAt = sharedPreferences.getString(KEY_USER_CREATED_AT, "") ?: ""
            )
        }
    }

    fun getAuthToken(): String? {
        return sharedPreferences.getString(KEY_JWT_TOKEN, null)
    }

    fun saveUserSession(token: String, user: User) {
        sharedPreferences.edit()
            .putString(KEY_JWT_TOKEN, token)
            .putString(KEY_USER_ID, user.id)
            .putString(KEY_USER_EMAIL, user.email)
            .putString(KEY_USER_CREATED_AT, user.createdAt)
            .apply()
        
        _currentUser.value = user
        _authState.value = AuthState.Authenticated
    }

    fun clearSession() {
        sharedPreferences.edit().clear().apply()
        _currentUser.value = null
        _authState.value = AuthState.Unauthenticated
    }

    fun setLoading() {
        _authState.value = AuthState.Loading
    }

    fun setError(message: String) {
        _authState.value = AuthState.Error(message)
    }

    private fun hasValidSession(): Boolean {
        val token = sharedPreferences.getString(KEY_JWT_TOKEN, null)
        val userId = sharedPreferences.getString(KEY_USER_ID, null)
        return !token.isNullOrEmpty() && !userId.isNullOrEmpty()
    }

    fun isAuthenticated(): Boolean {
        return _authState.value is AuthState.Authenticated
    }
}
