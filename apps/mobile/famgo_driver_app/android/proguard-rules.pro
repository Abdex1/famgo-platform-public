# Flutter
-keep class io.flutter.** { *; }

# Firebase
-keep class com.google.firebase.** { *; }
-dontwarn com.google.firebase.**

# Google Play Services
-keep class com.google.android.gms.** { *; }
-dontwarn com.google.android.gms.**

# Keep annotations
-keepattributes *Annotation*

# Retrofit / Gson (if used)
-keepattributes Signature
-keepattributes Exceptions