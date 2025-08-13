package com.germainleignel.days.ui.screens

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.text.input.PasswordVisualTransformation
import androidx.compose.ui.unit.dp
import androidx.lifecycle.viewmodel.compose.viewModel
import com.germainleignel.days.api.auth.UserSessionManager
import com.germainleignel.days.api.repository.ApiDataRepository
import com.germainleignel.days.storage.RepositoryFactory
import com.germainleignel.days.viewmodel.AuthViewModel
import kotlinx.coroutines.launch

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun AuthScreen(
    onAuthSuccess: () -> Unit,
    modifier: Modifier = Modifier
) {
    val context = LocalContext.current
    val sessionManager = RepositoryFactory.getSessionManager(context)
    val authViewModel: AuthViewModel = viewModel(
        factory = AuthViewModel.Factory(sessionManager, context)
    )
    
    val authState by sessionManager.authState.collectAsState()
    val coroutineScope = rememberCoroutineScope()
    
    var email by remember { mutableStateOf("") }
    var password by remember { mutableStateOf("") }
    var isRegistering by remember { mutableStateOf(false) }
    
    // Navigate on successful authentication
    LaunchedEffect(authState) {
        if (authState is UserSessionManager.AuthState.Authenticated) {
            onAuthSuccess()
        }
    }
    
    Column(
        modifier = modifier
            .fillMaxSize()
            .padding(16.dp),
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.Center
    ) {
        Text(
            text = if (isRegistering) "Create Account" else "Sign In",
            style = MaterialTheme.typography.headlineMedium,
            modifier = Modifier.padding(bottom = 32.dp)
        )
        
        OutlinedTextField(
            value = email,
            onValueChange = { email = it },
            label = { Text("Email") },
            keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Email),
            modifier = Modifier.fillMaxWidth()
        )
        
        Spacer(modifier = Modifier.height(16.dp))
        
        OutlinedTextField(
            value = password,
            onValueChange = { password = it },
            label = { Text("Password") },
            visualTransformation = PasswordVisualTransformation(),
            keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Password),
            modifier = Modifier.fillMaxWidth()
        )
        
        Spacer(modifier = Modifier.height(24.dp))
        
        Button(
            onClick = {
                coroutineScope.launch {
                    if (isRegistering) {
                        authViewModel.register(email, password)
                    } else {
                        authViewModel.login(email, password)
                    }
                }
            },
            enabled = email.isNotBlank() && password.isNotBlank() && 
                    authState !is UserSessionManager.AuthState.Loading,
            modifier = Modifier.fillMaxWidth()
        ) {
            if (authState is UserSessionManager.AuthState.Loading) {
                CircularProgressIndicator(
                    modifier = Modifier.size(16.dp),
                    color = MaterialTheme.colorScheme.onPrimary
                )
            } else {
                Text(if (isRegistering) "Create Account" else "Sign In")
            }
        }
        
        Spacer(modifier = Modifier.height(16.dp))
        
        TextButton(
            onClick = { isRegistering = !isRegistering }
        ) {
            Text(
                if (isRegistering) "Already have an account? Sign in" 
                else "Don't have an account? Create one"
            )
        }
        
        // Error display
        if (authState is UserSessionManager.AuthState.Error) {
            Spacer(modifier = Modifier.height(16.dp))
            Card(
                colors = CardDefaults.cardColors(
                    containerColor = MaterialTheme.colorScheme.errorContainer
                )
            ) {
                Text(
                    text = (authState as UserSessionManager.AuthState.Error).message,
                    color = MaterialTheme.colorScheme.onErrorContainer,
                    modifier = Modifier.padding(16.dp)
                )
            }
        }
        
        Spacer(modifier = Modifier.height(32.dp))
        
        // Switch to local mode option
        TextButton(
            onClick = {
                RepositoryFactory.switchToLocalMode(context)
                onAuthSuccess()
            }
        ) {
            Text("Continue with Local Storage")
        }
    }
}
