PS C:\Users\FEMOS\Desktop\Femos\extrac\famGo\famgo_driver_app\android> .\gradlew signingReport

> Configure project :gradle
WARNING: Unsupported Kotlin plugin version.
The `embedded-kotlin` and `kotlin-dsl` plugins rely on features of Kotlin `2.2.0` that might work differently than in the requested version `2.2.20`.

> Configure project :app
Warning: Flutter support for your project's Android Gradle Plugin version (Android Gradle Plugin version 8.7.3) will soon be dropped. Please upgrade your Android Gradle Plugin version to a version of at least Android Gradle Plugin version 8.11.1 soon.
Alternatively, use the flag "--android-skip-build-dependency-validation" to bypass this check.

Potential fix: Your project's AGP version is typically defined in the plugins block of the `settings.gradle` file (C:\Users\FEMOS\Desktop\Femos\extrac\famGo\famgo_driver_app\android/settings.gradle), by a plugin with the id of com.android.application.
If you don't see a plugins block, your project was likely created with an older template version. In this case it is most likely defined in the top-level build.gradle file (C:\Users\FEMOS\Desktop\Femos\extrac\famGo\famgo_driver_app\android/build.gradle) by the following line in the dependencies block of the buildscript: "classpath 'com.android.tools.build:gradle:<version>'".

WARNING: Your Android app project: app located at: C:\Users\FEMOS\Desktop\Femos\extrac\famGo\famgo_driver_app\android\app\build.gradle.kts
applies the Kotlin Gradle Plugin, which will cause build failures in future versions of Flutter.
Please migrate your app to Built-in Kotlin using this guide: https://docs.flutter.dev/release/breaking-changes/migrate-to-built-in-kotlin/for-app-developers

WARNING: Your app uses the following plugins that apply Kotlin Gradle Plugin (KGP): firebase_database, firebase_storage, package_info_plus, restart_app
Future versions of Flutter will fail to build if your app uses plugins that apply KGP.

Please check the changelogs of these plugins and upgrade to a version that supports Built-in Kotlin.
If no such version exists, report the issue to the plugin. If necessary, here is a guide on filing
an issue against a plugin: https://docs.flutter.dev/release/breaking-changes/migrate-to-built-in-kotlin/for-app-developers#report-incompatible-kotlin-gradle-plugin-usage-to-plugin-authors

If you are a plugin author, please migrate your plugin to Built-in Kotlin using this guide: https://docs.flutter.dev/release/breaking-changes/migrate-to-built-in-kotlin/for-plugin-authors

> Task :app:signingReport
Variant: debug
Config: debug
Store: C:\Users\FEMOS\.android\debug.keystore
Alias: AndroidDebugKey
MD5: 8A:3B:0D:36:BF:18:42:36:96:51:77:2A:94:88:DB:0C
SHA1: 7A:5F:D8:4A:48:07:D6:94:6F:E6:26:F2:BB:6C:87:79:FC:18:9E:07
SHA-256: AB:F5:0C:7E:34:BF:9A:7D:1C:43:EE:B9:83:1A:F4:E2:6D:8E:B7:8A:4E:D9:4B:0D:78:DD:18:1C:AE:72:BF:23
Valid until: Tuesday, May 30, 2056
----------
Variant: release
Config: release
Store: C:\Users\FEMOS\Desktop\Femos\extrac\famGo\famgo_driver_app\android\keystore\famgo-release.jks
Alias: famgo
MD5: E6:48:BF:53:94:7B:0C:4D:F7:81:5E:E1:24:36:A0:3D
SHA1: 35:04:7C:C0:E7:0D:58:A2:BD:8A:A3:21:8D:10:BD:57:10:84:5D:81
SHA-256: 30:14:18:4B:6C:C8:28:FD:0B:76:20:B8:9C:55:90:35:72:3D:C6:EE:B4:F4:1F:9C:CF:BF:E2:17:39:82:C9:FE
Valid until: Sunday, November 2, 2053
----------
Variant: profile
Config: debug
Store: C:\Users\FEMOS\.android\debug.keystore
Alias: AndroidDebugKey
MD5: 8A:3B:0D:36:BF:18:42:36:96:51:77:2A:94:88:DB:0C
SHA1: 7A:5F:D8:4A:48:07:D6:94:6F:E6:26:F2:BB:6C:87:79:FC:18:9E:07
SHA-256: AB:F5:0C:7E:34:BF:9A:7D:1C:43:EE:B9:83:1A:F4:E2:6D:8E:B7:8A:4E:D9:4B:0D:78:DD:18:1C:AE:72:BF:23
Valid until: Tuesday, May 30, 2056
----------
Variant: debugAndroidTest
Config: debug
Store: C:\Users\FEMOS\.android\debug.keystore
Alias: AndroidDebugKey
MD5: 8A:3B:0D:36:BF:18:42:36:96:51:77:2A:94:88:DB:0C
SHA1: 7A:5F:D8:4A:48:07:D6:94:6F:E6:26:F2:BB:6C:87:79:FC:18:9E:07
SHA-256: AB:F5:0C:7E:34:BF:9A:7D:1C:43:EE:B9:83:1A:F4:E2:6D:8E:B7:8A:4E:D9:4B:0D:78:DD:18:1C:AE:72:BF:23
Valid until: Tuesday, May 30, 2056
----------

> Task :cloud_firestore:signingReport
Variant: debugAndroidTest
Config: debug
Store: C:\Users\FEMOS\.android\debug.keystore
Alias: AndroidDebugKey
MD5: 8A:3B:0D:36:BF:18:42:36:96:51:77:2A:94:88:DB:0C
SHA1: 7A:5F:D8:4A:48:07:D6:94:6F:E6:26:F2:BB:6C:87:79:FC:18:9E:07
SHA-256: AB:F5:0C:7E:34:BF:9A:7D:1C:43:EE:B9:83:1A:F4:E2:6D:8E:B7:8A:4E:D9:4B:0D:78:DD:18:1C:AE:72:BF:23
Valid until: Tuesday, May 30, 2056
----------

> Task :connectivity_plus:signingReport
Variant: debugAndroidTest
Config: debug
Store: C:\Users\FEMOS\.android\debug.keystore
Alias: AndroidDebugKey
MD5: 8A:3B:0D:36:BF:18:42:36:96:51:77:2A:94:88:DB:0C
SHA1: 7A:5F:D8:4A:48:07:D6:94:6F:E6:26:F2:BB:6C:87:79:FC:18:9E:07
SHA-256: AB:F5:0C:7E:34:BF:9A:7D:1C:43:EE:B9:83:1A:F4:E2:6D:8E:B7:8A:4E:D9:4B:0D:78:DD:18:1C:AE:72:BF:23
Valid until: Tuesday, May 30, 2056
----------

> Task :firebase_auth:signingReport
Variant: debugAndroidTest
Config: debug
Store: C:\Users\FEMOS\.android\debug.keystore
Alias: AndroidDebugKey
MD5: 8A:3B:0D:36:BF:18:42:36:96:51:77:2A:94:88:DB:0C
SHA1: 7A:5F:D8:4A:48:07:D6:94:6F:E6:26:F2:BB:6C:87:79:FC:18:9E:07
SHA-256: AB:F5:0C:7E:34:BF:9A:7D:1C:43:EE:B9:83:1A:F4:E2:6D:8E:B7:8A:4E:D9:4B:0D:78:DD:18:1C:AE:72:BF:23
Valid until: Tuesday, May 30, 2056
----------

> Task :firebase_core:signingReport
Variant: debugAndroidTest
Config: debug
Store: C:\Users\FEMOS\.android\debug.keystore
Alias: AndroidDebugKey
MD5: 8A:3B:0D:36:BF:18:42:36:96:51:77:2A:94:88:DB:0C
SHA1: 7A:5F:D8:4A:48:07:D6:94:6F:E6:26:F2:BB:6C:87:79:FC:18:9E:07
SHA-256: AB:F5:0C:7E:34:BF:9A:7D:1C:43:EE:B9:83:1A:F4:E2:6D:8E:B7:8A:4E:D9:4B:0D:78:DD:18:1C:AE:72:BF:23
Valid until: Tuesday, May 30, 2056
----------

> Task :firebase_database:signingReport
Variant: debugAndroidTest
Config: debug
Store: C:\Users\FEMOS\.android\debug.keystore
Alias: AndroidDebugKey
MD5: 8A:3B:0D:36:BF:18:42:36:96:51:77:2A:94:88:DB:0C
SHA1: 7A:5F:D8:4A:48:07:D6:94:6F:E6:26:F2:BB:6C:87:79:FC:18:9E:07
SHA-256: AB:F5:0C:7E:34:BF:9A:7D:1C:43:EE:B9:83:1A:F4:E2:6D:8E:B7:8A:4E:D9:4B:0D:78:DD:18:1C:AE:72:BF:23
Valid until: Tuesday, May 30, 2056
----------

> Task :firebase_messaging:signingReport
Variant: debugAndroidTest
Config: debug
Store: C:\Users\FEMOS\.android\debug.keystore
Alias: AndroidDebugKey
MD5: 8A:3B:0D:36:BF:18:42:36:96:51:77:2A:94:88:DB:0C
SHA1: 7A:5F:D8:4A:48:07:D6:94:6F:E6:26:F2:BB:6C:87:79:FC:18:9E:07
SHA-256: AB:F5:0C:7E:34:BF:9A:7D:1C:43:EE:B9:83:1A:F4:E2:6D:8E:B7:8A:4E:D9:4B:0D:78:DD:18:1C:AE:72:BF:23
Valid until: Tuesday, May 30, 2056
----------

> Task :firebase_storage:signingReport
Variant: debugAndroidTest
Config: debug
Store: C:\Users\FEMOS\.android\debug.keystore
Alias: AndroidDebugKey
MD5: 8A:3B:0D:36:BF:18:42:36:96:51:77:2A:94:88:DB:0C
SHA1: 7A:5F:D8:4A:48:07:D6:94:6F:E6:26:F2:BB:6C:87:79:FC:18:9E:07
SHA-256: AB:F5:0C:7E:34:BF:9A:7D:1C:43:EE:B9:83:1A:F4:E2:6D:8E:B7:8A:4E:D9:4B:0D:78:DD:18:1C:AE:72:BF:23
Valid until: Tuesday, May 30, 2056
----------

> Task :flutter_geofire:signingReport
Variant: debugAndroidTest
Config: debug
Store: C:\Users\FEMOS\.android\debug.keystore
Alias: AndroidDebugKey
MD5: 8A:3B:0D:36:BF:18:42:36:96:51:77:2A:94:88:DB:0C
SHA1: 7A:5F:D8:4A:48:07:D6:94:6F:E6:26:F2:BB:6C:87:79:FC:18:9E:07
SHA-256: AB:F5:0C:7E:34:BF:9A:7D:1C:43:EE:B9:83:1A:F4:E2:6D:8E:B7:8A:4E:D9:4B:0D:78:DD:18:1C:AE:72:BF:23
Valid until: Tuesday, May 30, 2056
----------

> Task :flutter_notification_channel:signingReport
Variant: debugAndroidTest
Config: debug
Store: C:\Users\FEMOS\.android\debug.keystore
Alias: AndroidDebugKey
MD5: 8A:3B:0D:36:BF:18:42:36:96:51:77:2A:94:88:DB:0C
SHA1: 7A:5F:D8:4A:48:07:D6:94:6F:E6:26:F2:BB:6C:87:79:FC:18:9E:07
SHA-256: AB:F5:0C:7E:34:BF:9A:7D:1C:43:EE:B9:83:1A:F4:E2:6D:8E:B7:8A:4E:D9:4B:0D:78:DD:18:1C:AE:72:BF:23
Valid until: Tuesday, May 30, 2056
----------

> Task :flutter_plugin_android_lifecycle:signingReport
Variant: debugAndroidTest
Config: debug
Store: C:\Users\FEMOS\.android\debug.keystore
Alias: AndroidDebugKey
MD5: 8A:3B:0D:36:BF:18:42:36:96:51:77:2A:94:88:DB:0C
SHA1: 7A:5F:D8:4A:48:07:D6:94:6F:E6:26:F2:BB:6C:87:79:FC:18:9E:07
SHA-256: AB:F5:0C:7E:34:BF:9A:7D:1C:43:EE:B9:83:1A:F4:E2:6D:8E:B7:8A:4E:D9:4B:0D:78:DD:18:1C:AE:72:BF:23
Valid until: Tuesday, May 30, 2056
----------

> Task :geolocator_android:signingReport
Variant: debugAndroidTest
Config: debug
Store: C:\Users\FEMOS\.android\debug.keystore
Alias: AndroidDebugKey
MD5: 8A:3B:0D:36:BF:18:42:36:96:51:77:2A:94:88:DB:0C
SHA1: 7A:5F:D8:4A:48:07:D6:94:6F:E6:26:F2:BB:6C:87:79:FC:18:9E:07
SHA-256: AB:F5:0C:7E:34:BF:9A:7D:1C:43:EE:B9:83:1A:F4:E2:6D:8E:B7:8A:4E:D9:4B:0D:78:DD:18:1C:AE:72:BF:23
Valid until: Tuesday, May 30, 2056
----------

> Task :google_maps_flutter_android:signingReport
Variant: debugAndroidTest
Config: debug
Store: C:\Users\FEMOS\.android\debug.keystore
Alias: AndroidDebugKey
MD5: 8A:3B:0D:36:BF:18:42:36:96:51:77:2A:94:88:DB:0C
SHA1: 7A:5F:D8:4A:48:07:D6:94:6F:E6:26:F2:BB:6C:87:79:FC:18:9E:07
SHA-256: AB:F5:0C:7E:34:BF:9A:7D:1C:43:EE:B9:83:1A:F4:E2:6D:8E:B7:8A:4E:D9:4B:0D:78:DD:18:1C:AE:72:BF:23
Valid until: Tuesday, May 30, 2056
----------

> Task :google_sign_in_android:signingReport
Variant: debugAndroidTest
Config: debug
Store: C:\Users\FEMOS\.android\debug.keystore
Alias: AndroidDebugKey
MD5: 8A:3B:0D:36:BF:18:42:36:96:51:77:2A:94:88:DB:0C
SHA1: 7A:5F:D8:4A:48:07:D6:94:6F:E6:26:F2:BB:6C:87:79:FC:18:9E:07
SHA-256: AB:F5:0C:7E:34:BF:9A:7D:1C:43:EE:B9:83:1A:F4:E2:6D:8E:B7:8A:4E:D9:4B:0D:78:DD:18:1C:AE:72:BF:23
Valid until: Tuesday, May 30, 2056
----------

> Task :image_cropper:signingReport
Variant: debugAndroidTest
Config: debug
Store: C:\Users\FEMOS\.android\debug.keystore
Alias: AndroidDebugKey
MD5: 8A:3B:0D:36:BF:18:42:36:96:51:77:2A:94:88:DB:0C
SHA1: 7A:5F:D8:4A:48:07:D6:94:6F:E6:26:F2:BB:6C:87:79:FC:18:9E:07
SHA-256: AB:F5:0C:7E:34:BF:9A:7D:1C:43:EE:B9:83:1A:F4:E2:6D:8E:B7:8A:4E:D9:4B:0D:78:DD:18:1C:AE:72:BF:23
Valid until: Tuesday, May 30, 2056
----------

> Task :image_picker_android:signingReport
Variant: debugAndroidTest
Config: debug
Store: C:\Users\FEMOS\.android\debug.keystore
Alias: AndroidDebugKey
MD5: 8A:3B:0D:36:BF:18:42:36:96:51:77:2A:94:88:DB:0C
SHA1: 7A:5F:D8:4A:48:07:D6:94:6F:E6:26:F2:BB:6C:87:79:FC:18:9E:07
SHA-256: AB:F5:0C:7E:34:BF:9A:7D:1C:43:EE:B9:83:1A:F4:E2:6D:8E:B7:8A:4E:D9:4B:0D:78:DD:18:1C:AE:72:BF:23
Valid until: Tuesday, May 30, 2056
----------

> Task :jni:signingReport
Variant: debugAndroidTest
Config: debug
Store: C:\Users\FEMOS\.android\debug.keystore
Alias: AndroidDebugKey
MD5: 8A:3B:0D:36:BF:18:42:36:96:51:77:2A:94:88:DB:0C
SHA1: 7A:5F:D8:4A:48:07:D6:94:6F:E6:26:F2:BB:6C:87:79:FC:18:9E:07
SHA-256: AB:F5:0C:7E:34:BF:9A:7D:1C:43:EE:B9:83:1A:F4:E2:6D:8E:B7:8A:4E:D9:4B:0D:78:DD:18:1C:AE:72:BF:23
Valid until: Tuesday, May 30, 2056
----------

> Task :jni_flutter:signingReport
Variant: debugAndroidTest
Config: debug
Store: C:\Users\FEMOS\.android\debug.keystore
Alias: AndroidDebugKey
MD5: 8A:3B:0D:36:BF:18:42:36:96:51:77:2A:94:88:DB:0C
SHA1: 7A:5F:D8:4A:48:07:D6:94:6F:E6:26:F2:BB:6C:87:79:FC:18:9E:07
SHA-256: AB:F5:0C:7E:34:BF:9A:7D:1C:43:EE:B9:83:1A:F4:E2:6D:8E:B7:8A:4E:D9:4B:0D:78:DD:18:1C:AE:72:BF:23
Valid until: Tuesday, May 30, 2056
----------

> Task :package_info_plus:signingReport
Variant: debugAndroidTest
Config: debug
Store: C:\Users\FEMOS\.android\debug.keystore
Alias: AndroidDebugKey
MD5: 8A:3B:0D:36:BF:18:42:36:96:51:77:2A:94:88:DB:0C
SHA1: 7A:5F:D8:4A:48:07:D6:94:6F:E6:26:F2:BB:6C:87:79:FC:18:9E:07
SHA-256: AB:F5:0C:7E:34:BF:9A:7D:1C:43:EE:B9:83:1A:F4:E2:6D:8E:B7:8A:4E:D9:4B:0D:78:DD:18:1C:AE:72:BF:23
Valid until: Tuesday, May 30, 2056
----------

> Task :permission_handler_android:signingReport
Variant: debugAndroidTest
Config: debug
Store: C:\Users\FEMOS\.android\debug.keystore
Alias: AndroidDebugKey
MD5: 8A:3B:0D:36:BF:18:42:36:96:51:77:2A:94:88:DB:0C
SHA1: 7A:5F:D8:4A:48:07:D6:94:6F:E6:26:F2:BB:6C:87:79:FC:18:9E:07
SHA-256: AB:F5:0C:7E:34:BF:9A:7D:1C:43:EE:B9:83:1A:F4:E2:6D:8E:B7:8A:4E:D9:4B:0D:78:DD:18:1C:AE:72:BF:23
Valid until: Tuesday, May 30, 2056
----------

> Task :restart_app:signingReport
Variant: debugAndroidTest
Config: debug
Store: C:\Users\FEMOS\.android\debug.keystore
Alias: AndroidDebugKey
MD5: 8A:3B:0D:36:BF:18:42:36:96:51:77:2A:94:88:DB:0C
SHA1: 7A:5F:D8:4A:48:07:D6:94:6F:E6:26:F2:BB:6C:87:79:FC:18:9E:07
SHA-256: AB:F5:0C:7E:34:BF:9A:7D:1C:43:EE:B9:83:1A:F4:E2:6D:8E:B7:8A:4E:D9:4B:0D:78:DD:18:1C:AE:72:BF:23
Valid until: Tuesday, May 30, 2056
----------

> Task :shared_preferences_android:signingReport
Variant: debugAndroidTest
Config: debug
Store: C:\Users\FEMOS\.android\debug.keystore
Alias: AndroidDebugKey
MD5: 8A:3B:0D:36:BF:18:42:36:96:51:77:2A:94:88:DB:0C
SHA1: 7A:5F:D8:4A:48:07:D6:94:6F:E6:26:F2:BB:6C:87:79:FC:18:9E:07
SHA-256: AB:F5:0C:7E:34:BF:9A:7D:1C:43:EE:B9:83:1A:F4:E2:6D:8E:B7:8A:4E:D9:4B:0D:78:DD:18:1C:AE:72:BF:23
Valid until: Tuesday, May 30, 2056
----------

> Task :sqflite_android:signingReport
Variant: debugAndroidTest
Config: debug
Store: C:\Users\FEMOS\.android\debug.keystore
Alias: AndroidDebugKey
MD5: 8A:3B:0D:36:BF:18:42:36:96:51:77:2A:94:88:DB:0C
SHA1: 7A:5F:D8:4A:48:07:D6:94:6F:E6:26:F2:BB:6C:87:79:FC:18:9E:07
SHA-256: AB:F5:0C:7E:34:BF:9A:7D:1C:43:EE:B9:83:1A:F4:E2:6D:8E:B7:8A:4E:D9:4B:0D:78:DD:18:1C:AE:72:BF:23
Valid until: Tuesday, May 30, 2056
----------

> Task :url_launcher_android:signingReport
Variant: debugAndroidTest
Config: debug
Store: C:\Users\FEMOS\.android\debug.keystore
Alias: AndroidDebugKey
MD5: 8A:3B:0D:36:BF:18:42:36:96:51:77:2A:94:88:DB:0C
SHA1: 7A:5F:D8:4A:48:07:D6:94:6F:E6:26:F2:BB:6C:87:79:FC:18:9E:07
SHA-256: AB:F5:0C:7E:34:BF:9A:7D:1C:43:EE:B9:83:1A:F4:E2:6D:8E:B7:8A:4E:D9:4B:0D:78:DD:18:1C:AE:72:BF:23
Valid until: Tuesday, May 30, 2056
----------

[Incubating] Problems report is available at: file:///C:/Users/FEMOS/Desktop/Femos/extrac/famGo/famgo_driver_app/build/reports/problems/problems-report.html

Deprecated Gradle features were used in this build, making it incompatible with Gradle 10.

You can use '--warning-mode all' to show the individual deprecation warnings and determine if they come from your own scripts or plugins.

For more on this, please refer to https://docs.gradle.org/9.1.0/userguide/command_line_interface.html#sec:command_line_warnings in the Gradle documentation.

BUILD SUCCESSFUL in 53s
28 actionable tasks: 24 executed, 4 up-to-date
Consider enabling configuration cache to speed up this build: https://docs.gradle.org/9.1.0/userguide/configuration_cache_enabling.html