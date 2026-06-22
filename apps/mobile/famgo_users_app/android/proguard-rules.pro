# Firebase
-keepattributes Signature
-keepattributes *Annotation*

# Google Sign In
-keep class com.google.android.gms.auth.api.signin.** { *; }

# Kotlin
-keep class kotlin.Metadata { *; }

# Stripe
-dontwarn com.stripe.**
-keep class com.stripe.** { *; }

-dontwarn com.stripe.android.pushProvisioning.PushProvisioningActivity$g
-dontwarn com.stripe.android.pushProvisioning.PushProvisioningActivityStarter$Args
-dontwarn com.stripe.android.pushProvisioning.PushProvisioningActivityStarter$Error
-dontwarn com.stripe.android.pushProvisioning.PushProvisioningActivityStarter
-dontwarn com.stripe.android.pushProvisioning.PushProvisioningEphemeralKeyProvider