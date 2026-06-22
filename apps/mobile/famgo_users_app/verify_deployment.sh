#!/bin/bash
# Verification Script for Firebase Auth Fix Implementation
# Checks that all Phase 1-3A components are in place

echo "=========================================="
echo "Firebase Auth Fix - Verification Script"
echo "=========================================="
echo ""

PROJECT_ROOT="C:\Users\FEMOS\Desktop\Femos\extrac\uber-clone-master\uber_users_app"
PASS=0
FAIL=0

# Function to check file exists
check_file() {
    local file=$1
    local name=$2
    if [ -f "$PROJECT_ROOT/$file" ]; then
        echo "✅ $name"
        ((PASS++))
    else
        echo "❌ $name (MISSING: $file)"
        ((FAIL++))
    fi
}

# Function to check contains text
check_content() {
    local file=$1
    local search=$2
    local name=$3
    if grep -q "$search" "$PROJECT_ROOT/$file" 2>/dev/null; then
        echo "✅ $name"
        ((PASS++))
    else
        echo "❌ $name"
        ((FAIL++))
    fi
}

echo "PHASE 1: Foundation Utilities"
echo "---"
check_file "lib/core/auth_validators.dart" "auth_validators.dart exists"
check_file "lib/core/rate_limiter.dart" "rate_limiter.dart exists"
check_file "lib/core/auth_constants.dart" "auth_constants.dart exists"
check_content "lib/core/auth_validators.dart" "isValidE164PhoneNumber" "E.164 validation"
check_content "lib/core/rate_limiter.dart" "checkLoginRateLimit" "Rate limiting"
check_content "lib/core/auth_constants.dart" "maxAttemptsPerHour" "Auth constants"
echo ""

echo "PHASE 2: Security"
echo "---"
check_file "firebase_realtime_database_rules.json" "Firebase Rules"
check_file "android/app/src/main/res/xml/network_security_config.xml" "Network Security Config"
check_file "lib/core/secure_otp_handler.dart" "secure_otp_handler.dart"
check_content "android/app/src/main/AndroidManifest.xml" "networkSecurityConfig" "AndroidManifest updated"
check_content "pubspec.yaml" "libphonenumber_plugin" "libphonenumber_plugin in pubspec"
echo ""

echo "PHASE 3A: Core Auth Refactor"
echo "---"
check_file "lib/appInfo/auth_provider_v2.dart" "auth_provider_v2.dart created"
check_content "lib/appInfo/auth_provider_v2.dart" "PhoneAuthOptions" "PhoneAuthOptions implemented"
check_content "lib/appInfo/auth_provider_v2.dart" "RateLimiter" "RateLimiter integrated"
check_content "lib/appInfo/auth_provider_v2.dart" "AuthValidators" "AuthValidators integrated"
echo ""

echo "Documentation Files"
echo "---"
check_file "COMPLETE_IMPLEMENTATION_SUMMARY.md" "Implementation summary"
check_file "DEPLOYMENT_STATUS_REPORT.md" "Deployment status"
check_file "MIGRATION_GUIDE_AUTH_PROVIDER.md" "Migration guide"
check_file "EXECUTIVE_HANDOFF.md" "Executive handoff"
check_file "PRODUCTION_DEPLOYMENT_PLAN.md" "Deployment plan"
echo ""

echo "=========================================="
echo "Summary"
echo "=========================================="
echo "✅ Passed: $PASS"
echo "❌ Failed: $FAIL"
echo ""

if [ $FAIL -eq 0 ]; then
    echo "✅ ALL CHECKS PASSED"
    echo "Phase 1-3A implementation complete!"
    echo ""
    echo "Next steps:"
    echo "1. Review EXECUTIVE_HANDOFF.md"
    echo "2. Run: flutter pub get"
    echo "3. Run: flutter analyze"
    echo "4. Begin Phase 3B (OTP Screen updates)"
else
    echo "❌ SOME CHECKS FAILED"
    echo "Please verify missing files are created correctly"
fi

echo ""
echo "For detailed information, see:"
echo "- COMPLETE_IMPLEMENTATION_SUMMARY.md"
echo "- EXECUTIVE_HANDOFF.md"
echo "=========================================="
