import java.util.Properties
val keystoreProperties = Properties()

val keystorePropertiesFile =
    rootProject.file("key.properties")

if (keystorePropertiesFile.exists()) {
    keystoreProperties.load(
        keystorePropertiesFile.inputStream()
    )
}
plugins {
    id("com.android.application")
    id("org.jetbrains.kotlin.android")
    id("com.google.gms.google-services")
    id("dev.flutter.flutter-gradle-plugin")
}

android {
    namespace = "com.famgo.famgo_drivers_app"
    compileSdk = 36
    ndkVersion = flutter.ndkVersion

    compileOptions {
        sourceCompatibility = JavaVersion.VERSION_17
        targetCompatibility = JavaVersion.VERSION_17
    }

    defaultConfig {
        // TODO: Specify your own unique Application ID (https://developer.android.com/studio/build/application-id.html).
        applicationId = "com.famgo.famgo_drivers_app"
        // You can update the following values to match your application needs.
        // For more information, see: https://flutter.dev/to/review-gradle-config.
        minSdk = 26
        targetSdk = 36
        versionCode = flutter.versionCode
        versionName = flutter.versionName
    }

    signingConfigs {

    if (keystorePropertiesFile.exists()) {

        create("release") {

            storeFile = file(keystoreProperties["storeFile"] as String)
            storePassword = keystoreProperties["storePassword"] as String
            keyAlias = keystoreProperties["keyAlias"] as String
            keyPassword = keystoreProperties["keyPassword"] as String
        }
    }
}

buildTypes {

    release {

        signingConfig =
            if (keystorePropertiesFile.exists())
                signingConfigs.getByName("release")
            else
                signingConfigs.getByName("debug")

        isMinifyEnabled = true

        isShrinkResources = true

        proguardFiles(
            getDefaultProguardFile(
                "proguard-android-optimize.txt"
            ),
            "../proguard-rules.pro"
        )
    }

    debug {
        signingConfig = signingConfigs.getByName("debug")
    }
}
}
kotlin {
    compilerOptions {
        jvmTarget = org.jetbrains.kotlin.gradle.dsl.JvmTarget.JVM_17
    }
}

dependencies {

    implementation(platform("com.google.firebase:firebase-bom:34.14.1"))

    implementation("com.google.firebase:firebase-auth")

    implementation("com.google.firebase:firebase-database")

    implementation("com.google.firebase:firebase-firestore")

    implementation("com.google.firebase:firebase-storage")

    implementation("com.google.firebase:firebase-messaging")

    implementation("androidx.credentials:credentials:1.3.0")

    implementation("androidx.credentials:credentials-play-services-auth:1.3.0")

    implementation("com.google.android.libraries.identity.googleid:googleid:1.1.1")
    implementation("com.google.firebase:firebase-appcheck-debug")
}
flutter {
    source = "../.."
}
