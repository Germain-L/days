package com.germainleignel.days.storage

import android.content.Context

/**
 * Factory to create data repository instances.
 * This allows switching between local and remote storage implementations.
 */
object RepositoryFactory {

    @Volatile
    private var INSTANCE: DataRepository? = null

    /**
     * Get the current repository instance.
     * Defaults to local storage if no instance is set.
     */
    fun getRepository(context: Context): DataRepository {
        return INSTANCE ?: synchronized(this) {
            val instance = INSTANCE ?: createLocalRepository(context)
            INSTANCE = instance
            instance
        }
    }

    /**
     * Create a local storage repository
     */
    fun createLocalRepository(context: Context): DataRepository {
        return LocalDataRepository(context.applicationContext)
    }

    /**
     * Switch to a different repository implementation
     * (e.g., when backend becomes available)
     */
    fun switchRepository(newRepository: DataRepository) {
        synchronized(this) {
            INSTANCE = newRepository
        }
    }

    /**
     * Reset the repository instance (useful for testing)
     */
    fun resetRepository() {
        synchronized(this) {
            INSTANCE = null
        }
    }

    /**
     * Check if a repository is currently set
     */
    fun hasRepository(): Boolean = INSTANCE != null

    /**
     * Get the current storage type
     */
    fun getCurrentStorageType(context: Context): StorageType {
        return getRepository(context).storageType
    }
}

/**
 * Configuration for future backend integration
 */
data class BackendConfig(
    val baseUrl: String,
    val apiKey: String? = null,
    val enableOfflineMode: Boolean = true
)

/**
 * Interface for future remote repository implementation
 */
interface RemoteDataRepository : DataRepository {
    suspend fun syncWithLocal(localRepository: DataRepository): SyncResult
    suspend fun isOnline(): Boolean
    suspend fun getLastSyncTime(): Long?
}

/**
 * Result of sync operations between local and remote storage
 */
sealed class SyncResult {
    object Success : SyncResult()
    data class Failure(val error: String) : SyncResult()
    data class Conflict(val conflictedItems: List<String>) : SyncResult()
}
