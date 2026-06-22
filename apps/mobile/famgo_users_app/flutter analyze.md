
PS C:\Users\FEMOS\Desktop\Femos\extrac\uber-clone-master\uber_users_app> flutter analyze
Resolving dependencies...
Downloading packages...
  cli_util 0.4.2 (0.5.1 available)
  flutter_polyline_points 2.1.0 (3.1.0 available)
  google_sign_in 6.3.0 (7.2.0 available)
  google_sign_in_android 6.2.1 (7.2.13 available)
  google_sign_in_ios 5.9.0 (6.3.0 available)
  google_sign_in_platform_interface 2.5.0 (3.1.0 available)
  google_sign_in_web 0.12.4+4 (1.1.3 available)
  image 4.8.0 (4.9.1 available)
  js 0.6.7 (0.7.2 available)
  libphonenumber_platform_interface 0.3.1 (0.4.2 available)
  libphonenumber_plugin 0.2.5 (0.3.3 available)
  libphonenumber_web 0.2.0+1 (0.3.2 available)
  matcher 0.12.19 (0.12.20 available)
  meta 1.18.0 (1.18.3 available)
  package_config 2.2.0 (3.0.0 available)
  test_api 0.7.11 (0.7.12 available)
  vector_math 2.2.0 (2.4.0 available)
  xml 6.6.1 (7.0.1 available)
Got dependencies!
18 packages have newer versions incompatible with dependency constraints.
Try `flutter pub outdated` for more information.
Analyzing uber_users_app...

  error - The method 'PhoneAuthOptions' isn't defined for the type 'AuthenticationProvider'. Try
         correcting the name to the name of an existing method, or defining a method named
         'PhoneAuthOptions' - lib\appInfo\auth_provider.dart:99:32 - undefined_method
  error - The named parameter 'codeAutoRetrievalTimeout' is required, but there's no corresponding
         argument. Try adding the required argument - lib\appInfo\auth_provider.dart:152:35 -
         missing_required_argument
  error - The named parameter 'codeSent' is required, but there's no corresponding argument. Try
         adding the required argument - lib\appInfo\auth_provider.dart:152:35 -
         missing_required_argument
  error - The named parameter 'verificationCompleted' is required, but there's no corresponding
         argument. Try adding the required argument - lib\appInfo\auth_provider.dart:152:35 -
         missing_required_argument
  error - The named parameter 'verificationFailed' is required, but there's no corresponding
         argument. Try adding the required argument - lib\appInfo\auth_provider.dart:152:35 -
         missing_required_argument
  error - The named parameter 'phoneAuthOptions' isn't defined. Try correcting the name to an
         existing named parameter's name, or defining a named parameter with the name
         'phoneAuthOptions' - lib\appInfo\auth_provider.dart:153:9 - undefined_named_parameter
   info - Don't use 'BuildContext's across async gaps. Try rewriting the code to not use the
          'BuildContext', or guard the use with a 'mounted' check -
          lib\appInfo\auth_provider.dart:169:38 - use_build_context_synchronously
   info - 'onPopInvoked' is deprecated and shouldn't be used. Use onPopInvokedWithResult instead.
          This feature was deprecated after v3.22.0-12.0.pre. Try replacing the use of the deprecated
          member with the replacement - lib\authentication\otp_screen.dart:65:7 -
          deprecated_member_use
   info - Don't use 'BuildContext's across async gaps, guarded by an unrelated 'mounted' check. Guard
          a 'State.context' use with a 'mounted' check on the State, and other BuildContext use with
          a 'mounted' check on the BuildContext - lib\authentication\register_screen.dart:269:41 -
          use_build_context_synchronously
warning - Unused import: 'dart:io'. Try removing the import directive -
       lib\authentication\user_information_screen.dart:1:8 - unused_import
warning - Unused import: 'package:famgo_passenger_app/core/app_typography.dart'. Try removing the
       import directive - lib\authentication\user_information_screen.dart:9:8 - unused_import
warning - Unused import: 'package:famgo_passenger_app/core/app_shadows.dart'. Try removing the import
       directive - lib\authentication\user_information_screen.dart:10:8 - unused_import
   info - Parameter 'key' could be a super parameter. Trying converting 'key' to a super parameter -
          lib\authentication\user_information_screen.dart:14:9 - use_super_parameters
   info - Dangling library doc comment. Add a 'library' directive after the library comment -
          lib\core\auth_constants.dart:1:1 - dangling_library_doc_comments
   info - Dangling library doc comment. Add a 'library' directive after the library comment -
          lib\core\auth_validators.dart:1:1 - dangling_library_doc_comments
   info - Dangling library doc comment. Add a 'library' directive after the library comment -
          lib\core\rate_limiter.dart:1:1 - dangling_library_doc_comments
   info - Dangling library doc comment. Add a 'library' directive after the library comment -
          lib\core\secure_otp_handler.dart:1:1 - dangling_library_doc_comments
   info - Use an initializing formal to assign a parameter to a field. Try using an initialing formal
          ('this._code') to initialize the field - lib\core\secure_otp_handler.dart:70:9 -
          prefer_initializing_formals
   info - Use an initializing formal to assign a parameter to a field. Try using an initialing formal
          ('this._expirationDuration') to initialize the field -
          lib\core\secure_otp_handler.dart:72:9 - prefer_initializing_formals
warning - The value of the field '_verificationId' isn't used. Try removing the field, or using it -
       lib\core\secure_otp_handler.dart:123:15 - unused_field
warning - The value of the field '_resendAttempts' isn't used. Try removing the field, or using it -
       lib\core\secure_otp_handler.dart:125:7 - unused_field
   info - The imported package 'firebase_app_check' isn't a dependency of the importing package. Try
          adding a dependency for 'firebase_app_check' in the 'pubspec.yaml' file -
          lib\main.dart:144:8 - depend_on_referenced_packages
  error - Target of URI doesn't exist: 'package:firebase_app_check/firebase_app_check.dart'. Try
         creating the file referenced by the URI, or try using a URI for a file that does exist -
         lib\main.dart:144:8 - uri_does_not_exist
  error - Undefined name 'FirebaseAppCheck'. Try correcting the name to one that is defined, or
         defining the name - lib\main.dart:168:9 - undefined_identifier
  error - Undefined name 'AndroidProvider'. Try correcting the name to one that is defined, or
         defining the name - lib\main.dart:170:35 - undefined_identifier
  error - Undefined name 'AndroidProvider'. Try correcting the name to one that is defined, or
         defining the name - lib\main.dart:170:59 - undefined_identifier
  error - Undefined name 'AppleProvider'. Try correcting the name to one that is defined, or defining
         the name - lib\main.dart:172:33 - undefined_identifier
  error - Undefined name 'AppleProvider'. Try correcting the name to one that is defined, or defining
         the name - lib\main.dart:172:55 - undefined_identifier
   info - The import of 'package:flutter/foundation.dart' is unnecessary because all of the used
          elements are also provided by the import of 'package:flutter/material.dart'. Try removing
          the import directive - lib\methods\common_methods.dart:12:8 - unnecessary_import
   info - The import of 'package:flutter/foundation.dart' is unnecessary because all of the used
          elements are also provided by the import of 'package:flutter/material.dart'. Try removing
          the import directive - lib\methods\common_methods.dart:14:8 - unnecessary_import
warning - Duplicate import. Try removing all but one import of the library -
       lib\methods\common_methods.dart:14:8 - duplicate_import
   info - Missing type annotation. Try adding a type annotation -
          lib\methods\common_methods.dart:16:3 - strict_top_level_inference
   info - Missing type annotation. Try adding a type annotation -
          lib\methods\common_methods.dart:28:3 - strict_top_level_inference
   info - Missing type annotation. Try adding a type annotation -
          lib\methods\common_methods.dart:33:10 - strict_top_level_inference
   info - Don't use 'BuildContext's across async gaps. Try rewriting the code to not use the
          'BuildContext', or guard the use with a 'mounted' check -
          lib\methods\common_methods.dart:69:33 - use_build_context_synchronously
   info - Missing type annotation. Try adding a type annotation -
          lib\methods\common_methods.dart:131:3 - strict_top_level_inference
   info - The import of 'package:flutter/foundation.dart' is unnecessary because all of the used
          elements are also provided by the import of 'package:flutter/material.dart'. Try removing
          the import directive - lib\methods\push_notification_service.dart:2:8 - unnecessary_import
warning - Unused import: 'package:googleapis/servicecontrol/v1.dart'. Try removing the import
       directive - lib\methods\push_notification_service.dart:6:8 - unused_import
warning - Unused import: '../main.dart'. Try removing the import directive -
       lib\methods\push_notification_service.dart:10:8 - unused_import
   info - Missing type annotation. Try adding a type annotation -
          lib\methods\push_notification_service.dart:56:10 - strict_top_level_inference
   info - Unnecessary braces in a string interpolation. Try removing the braces -
          lib\methods\push_notification_service.dart:58:35 - unnecessary_brace_in_string_interps
   info - Unnecessary braces in a string interpolation. Try removing the braces -
          lib\methods\push_notification_service.dart:68:35 - unnecessary_brace_in_string_interps
   info - The variable name 'place_id' isn't a lowerCamelCase identifier. Try changing the name to
          follow the lowerCamelCase style - lib\models\prediction_model.dart:2:11 -
          non_constant_identifier_names
   info - The variable name 'main_text' isn't a lowerCamelCase identifier. Try changing the name to
          follow the lowerCamelCase style - lib\models\prediction_model.dart:3:11 -
          non_constant_identifier_names
   info - The variable name 'secondary_text' isn't a lowerCamelCase identifier. Try changing the name
          to follow the lowerCamelCase style - lib\models\prediction_model.dart:4:11 -
          non_constant_identifier_names
   info - The variable name 'place_id' isn't a lowerCamelCase identifier. Try changing the name to
          follow the lowerCamelCase style - lib\models\prediction_model.dart:6:25 -
          non_constant_identifier_names
   info - The variable name 'main_text' isn't a lowerCamelCase identifier. Try changing the name to
          follow the lowerCamelCase style - lib\models\prediction_model.dart:6:40 -
          non_constant_identifier_names
   info - The variable name 'secondary_text' isn't a lowerCamelCase identifier. Try changing the name
          to follow the lowerCamelCase style - lib\models\prediction_model.dart:6:56 -
          non_constant_identifier_names
warning - Unused import: 'package:famgo_passenger_app/core/app_colors.dart'. Try removing the import
       directive - lib\pages\about_page.dart:2:8 - unused_import
warning - Unused import: 'package:famgo_passenger_app/core/app_typography.dart'. Try removing the
       import directive - lib\pages\about_page.dart:3:8 - unused_import
warning - Unused import: 'package:famgo_passenger_app/core/app_shadows.dart'. Try removing the import
       directive - lib\pages\about_page.dart:4:8 - unused_import
warning - Unused import: '../main.dart'. Try removing the import directive -
       lib\pages\blocked_screen.dart:3:8 - unused_import
warning - Unused import: 'package:famgo_passenger_app/core/app_typography.dart'. Try removing the
       import directive - lib\pages\blocked_screen.dart:5:8 - unused_import
warning - Unused import: 'package:famgo_passenger_app/core/app_shadows.dart'. Try removing the import
       directive - lib\pages\blocked_screen.dart:6:8 - unused_import
   info - The import of 'dart:typed_data' is unnecessary because all of the used elements are also
          provided by the import of 'package:flutter/foundation.dart'. Try removing the import
          directive - lib\pages\home_page.dart:3:8 - unnecessary_import
warning - Unused import: 'package:famgo_passenger_app/pages/profile_page.dart'. Try removing the
       import directive - lib\pages\home_page.dart:19:8 - unused_import
warning - Unused import: 'package:famgo_passenger_app/widgets/sign_out_dialog.dart'. Try removing the
       import directive - lib\pages\home_page.dart:22:8 - unused_import
warning - Unused import: 'about_page.dart'. Try removing the import directive -
       lib\pages\home_page.dart:35:8 - unused_import
warning - Unused import: 'trips_history_page.dart'. Try removing the import directive -
       lib\pages\home_page.dart:36:8 - unused_import
warning - Unused import: '../main.dart'. Try removing the import directive -
       lib\pages\home_page.dart:37:8 - unused_import
warning - Unused import: 'package:famgo_passenger_app/core/app_typography.dart'. Try removing the
       import directive - lib\pages\home_page.dart:39:8 - unused_import
warning - Unused import: 'package:famgo_passenger_app/core/app_shadows.dart'. Try removing the import
       directive - lib\pages\home_page.dart:40:8 - unused_import
   info - Missing type annotation. Try adding a type annotation - lib\pages\home_page.dart:82:3 -
          strict_top_level_inference
warning - The value of the local variable 'configuration' isn't used. Try removing the variable or
       using it - lib\pages\home_page.dart:84:26 - unused_local_variable
   info - Missing type annotation. Try adding a type annotation - lib\pages\home_page.dart:104:3 -
          strict_top_level_inference
   info - 'setMapStyle' is deprecated and shouldn't be used. Use GoogleMap.style instead. Try
          replacing the use of the deprecated member with the replacement -
          lib\pages\home_page.dart:105:16 - deprecated_member_use
   info - Missing type annotation. Try adding a type annotation - lib\pages\home_page.dart:108:3 -
          strict_top_level_inference
   info - 'desiredAccuracy' is deprecated and shouldn't be used. use settings parameter with
          AndroidSettings, AppleSettings, WebSettings, or LocationSettings. Try replacing the use of
          the deprecated member with the replacement - lib\pages\home_page.dart:110:9 -
          deprecated_member_use
   info - Don't use 'BuildContext's across async gaps. Try rewriting the code to not use the
          'BuildContext', or guard the use with a 'mounted' check - lib\pages\home_page.dart:121:33 -
          use_build_context_synchronously
   info - Missing type annotation. Try adding a type annotation - lib\pages\home_page.dart:128:3 -
          strict_top_level_inference
   info - Don't use 'BuildContext's across async gaps. Try rewriting the code to not use the
          'BuildContext', or guard the use with a 'mounted' check - lib\pages\home_page.dart:147:26 -
          use_build_context_synchronously
   info - Don't use 'BuildContext's across async gaps. Try rewriting the code to not use the
          'BuildContext', or guard the use with a 'mounted' check - lib\pages\home_page.dart:151:70 -
          use_build_context_synchronously
   info - Don't use 'BuildContext's across async gaps. Try rewriting the code to not use the
          'BuildContext', or guard the use with a 'mounted' check - lib\pages\home_page.dart:156:13 -
          use_build_context_synchronously
   info - Missing type annotation. Try adding a type annotation - lib\pages\home_page.dart:161:3 -
          strict_top_level_inference
   info - Missing type annotation. Try adding a type annotation - lib\pages\home_page.dart:174:3 -
          strict_top_level_inference
   info - Don't use 'BuildContext's across async gaps. Try rewriting the code to not use the
          'BuildContext', or guard the use with a 'mounted' check - lib\pages\home_page.dart:204:19 -
          use_build_context_synchronously
   info - Function literals shouldn't be passed to 'forEach'. Try using a 'for' loop -
          lib\pages\home_page.dart:213:43 - avoid_function_literals_in_foreach_calls
   info - Missing type annotation. Try adding a type annotation - lib\pages\home_page.dart:323:3 -
          strict_top_level_inference
   info - Missing type annotation. Try adding a type annotation - lib\pages\home_page.dart:345:3 -
          strict_top_level_inference
   info - Missing type annotation. Try adding a type annotation - lib\pages\home_page.dart:355:3 -
          strict_top_level_inference
   info - Missing type annotation. Try adding a type annotation - lib\pages\home_page.dart:369:3 -
          strict_top_level_inference
   info - Unnecessary constructor invocation. Try using a collection literal -
          lib\pages\home_page.dart:376:34 - prefer_collection_literals
   info - Use interpolation to compose strings and values. Try using string interpolation to build
          the composite string - lib\pages\home_page.dart:385:13 -
          prefer_interpolation_to_compose_strings
   info - Missing type annotation. Try adding a type annotation - lib\pages\home_page.dart:400:3 -
          strict_top_level_inference
   info - Missing type annotation. Try adding a type annotation - lib\pages\home_page.dart:458:3 -
          strict_top_level_inference
   info - Unnecessary braces in a string interpolation. Try removing the braces -
          lib\pages\home_page.dart:511:31 - unnecessary_brace_in_string_interps
   info - Don't use 'BuildContext's across async gaps. Try rewriting the code to not use the
          'BuildContext', or guard the use with a 'mounted' check - lib\pages\home_page.dart:578:11 -
          use_build_context_synchronously
   info - Missing type annotation. Try adding a type annotation - lib\pages\home_page.dart:597:3 -
          strict_top_level_inference
   info - Missing type annotation. Try adding a type annotation - lib\pages\home_page.dart:605:3 -
          strict_top_level_inference
   info - Missing type annotation. Try adding a type annotation - lib\pages\home_page.dart:605:43 -
          strict_top_level_inference
   info - Missing type annotation. Try adding a type annotation - lib\pages\home_page.dart:639:3 -
          strict_top_level_inference
   info - Missing type annotation. Try adding a type annotation - lib\pages\home_page.dart:640:7 -
          strict_top_level_inference
   info - Missing type annotation. Try adding a type annotation - lib\pages\home_page.dart:677:3 -
          strict_top_level_inference
   info - Missing type annotation. Try adding a type annotation - lib\pages\home_page.dart:689:3 -
          strict_top_level_inference
   info - Use 'isEmpty' instead of 'length' to test whether the collection is empty. Try rewriting
          the expression to use 'isEmpty' - lib\pages\home_page.dart:690:9 - prefer_is_empty
   info - Missing type annotation. Try adding a type annotation - lib\pages\home_page.dart:705:3 -
          strict_top_level_inference
   info - Don't use 'BuildContext's across async gaps. Try rewriting the code to not use the
          'BuildContext', or guard the use with a 'mounted' check - lib\pages\home_page.dart:740:28 -
          use_build_context_synchronously
warning - Unused import: 'package:famgo_passenger_app/core/app_typography.dart'. Try removing the
       import directive - lib\pages\profile_page.dart:4:8 - unused_import
warning - Unused import: 'package:famgo_passenger_app/core/app_shadows.dart'. Try removing the import
       directive - lib\pages\profile_page.dart:5:8 - unused_import
   info - Missing type annotation. Try adding a type annotation - lib\pages\profile_page.dart:21:3 -
          strict_top_level_inference
   info - Parameter 'key' could be a super parameter. Trying converting 'key' to a super parameter -
          lib\pages\profile_widgets\profile_action_tile.dart:11:9 - use_super_parameters
   info - 'withOpacity' is deprecated and shouldn't be used. Use .withValues() to avoid precision
          loss. Try replacing the use of the deprecated member with the replacement -
          lib\pages\profile_widgets\profile_action_tile.dart:35:36 - deprecated_member_use
   info - 'withOpacity' is deprecated and shouldn't be used. Use .withValues() to avoid precision
          loss. Try replacing the use of the deprecated member with the replacement -
          lib\pages\profile_widgets\profile_action_tile.dart:36:45 - deprecated_member_use
   info - Parameter 'key' could be a super parameter. Trying converting 'key' to a super parameter -
          lib\pages\profile_widgets\profile_completion_card.dart:6:9 - use_super_parameters
   info - 'withOpacity' is deprecated and shouldn't be used. Use .withValues() to avoid precision
          loss. Try replacing the use of the deprecated member with the replacement -
          lib\pages\profile_widgets\profile_completion_card.dart:22:35 - deprecated_member_use
   info - Parameter 'key' could be a super parameter. Trying converting 'key' to a super parameter -
          lib\pages\profile_widgets\profile_hero_card.dart:13:9 - use_super_parameters
   info - 'withOpacity' is deprecated and shouldn't be used. Use .withValues() to avoid precision
          loss. Try replacing the use of the deprecated member with the replacement -
          lib\pages\profile_widgets\profile_hero_card.dart:168:35 - deprecated_member_use
   info - Use interpolation to compose strings and values. Try using string interpolation to build
          the composite string - lib\pages\profile_widgets\profile_hero_card.dart:186:31 -
          prefer_interpolation_to_compose_strings
   info - 'withOpacity' is deprecated and shouldn't be used. Use .withValues() to avoid precision
          loss. Try replacing the use of the deprecated member with the replacement -
          lib\pages\profile_widgets\profile_hero_card.dart:207:41 - deprecated_member_use
   info - Use interpolation to compose strings and values. Try using string interpolation to build
          the composite string - lib\pages\profile_widgets\profile_hero_card.dart:223:35 -
          prefer_interpolation_to_compose_strings
   info - The import of 'package:flutter/foundation.dart' is unnecessary because all of the used
          elements are also provided by the import of 'package:flutter/material.dart'. Try removing
          the import directive - lib\pages\search_destination_place.dart:9:8 - unnecessary_import
   info - Missing type annotation. Try adding a type annotation -
          lib\pages\search_destination_place.dart:23:3 - strict_top_level_inference
   info - Use interpolation to compose strings and values. Try using string interpolation to build
          the composite string - lib\pages\search_destination_place.dart:50:22 -
          prefer_interpolation_to_compose_strings
   info - Unnecessary braces in a string interpolation. Try removing the braces -
          lib\pages\search_destination_place.dart:63:39 - unnecessary_brace_in_string_interps
warning - Unused import: '../main.dart'. Try removing the import directive -
       lib\pages\trips_history_page.dart:4:8 - unused_import
   info - Use interpolation to compose strings and values. Try using string interpolation to build
          the composite string - lib\pages\trips_history_page.dart:118:33 -
          prefer_interpolation_to_compose_strings
   info - Unnecessary escape in string literal. Remove the '\' escape -
          lib\pages\trips_history_page.dart:118:34 - unnecessary_string_escapes
   info - Parameter 'key' could be a super parameter. Trying converting 'key' to a super parameter -
          lib\screens\splash_screen.dart:24:9 - use_super_parameters
   info - The import of 'package:flutter/foundation.dart' is unnecessary because all of the used
          elements are also provided by the import of 'package:flutter/material.dart'. Try removing
          the import directive - lib\services\stripe_payment_service.dart:6:8 - unnecessary_import
   info - Don't use 'BuildContext's across async gaps. Try rewriting the code to not use the
          'BuildContext', or guard the use with a 'mounted' check -
          lib\services\stripe_payment_service.dart:39:28 - use_build_context_synchronously
   info - Don't use 'BuildContext's across async gaps. Try rewriting the code to not use the
          'BuildContext', or guard the use with a 'mounted' check -
          lib\services\stripe_payment_service.dart:45:28 - use_build_context_synchronously
   info - Don't use 'BuildContext's across async gaps. Try rewriting the code to not use the
          'BuildContext', or guard the use with a 'mounted' check -
          lib\services\stripe_payment_service.dart:51:28 - use_build_context_synchronously
   info - Invalid use of a private type in a public API. Try making the private type public, or
          making the API that uses the private type also be private -
          lib\widgets\bid_dialog.dart:14:3 - library_private_types_in_public_api
warning - The value of the local variable 'fare' isn't used. Try removing the variable or using it -
       lib\widgets\bid_dialog.dart:91:12 - unused_local_variable
   info - Unnecessary braces in a string interpolation. Try removing the braces -
          lib\widgets\bid_dialog.dart:124:25 - unnecessary_brace_in_string_interps
warning - Unused import: 'package:googleapis/analytics/v3.dart'. Try removing the import directive -
       lib\widgets\custome_drawer.dart:2:8 - unused_import
warning - Unused import: '../main.dart'. Try removing the import directive -
       lib\widgets\custome_drawer.dart:9:8 - unused_import
   info - Parameter 'key' could be a super parameter. Trying converting 'key' to a super parameter -
          lib\widgets\custome_drawer.dart:15:9 - use_super_parameters
   info - Unnecessary empty statement. Try removing the empty statement or restructuring the code -
          lib\widgets\custome_drawer.dart:77:17 - empty_statements
   info - Unnecessary empty statement. Try removing the empty statement or restructuring the code -
          lib\widgets\custome_drawer.dart:83:17 - empty_statements
   info - Unnecessary empty statement. Try removing the empty statement or restructuring the code -
          lib\widgets\custome_drawer.dart:89:17 - empty_statements
   info - Unnecessary empty statement. Try removing the empty statement or restructuring the code -
          lib\widgets\custome_drawer.dart:112:17 - empty_statements
warning - This class (or a class that this class inherits from) is marked as '@immutable', but one or
       more of its instance fields aren't final: InfoDialog.title, InfoDialog.description -
       lib\widgets\info_dialog.dart:4:7 - must_be_immutable
   info - Parameter 'key' could be a super parameter. Trying converting 'key' to a super parameter -
          lib\widgets\loading_dialog.dart:6:9 - use_super_parameters
   info - Constructors in '@immutable' classes should be declared as 'const'. Try adding 'const' to
          the constructor declaration - lib\widgets\payment_dialog.dart:7:3 -
          prefer_const_constructors_in_immutables
   info - Unnecessary escape in string literal. Remove the '\' escape -
          lib\widgets\payment_dialog.dart:45:16 - unnecessary_string_escapes
   info - Don't use 'BuildContext's across async gaps. Try rewriting the code to not use the
          'BuildContext', or guard the use with a 'mounted' check -
          lib\widgets\payment_dialog.dart:106:11 - use_build_context_synchronously
   info - Don't use 'BuildContext's across async gaps. Try rewriting the code to not use the
          'BuildContext', or guard the use with a 'mounted' check -
          lib\widgets\payment_dialog.dart:113:13 - use_build_context_synchronously
   info - Don't invoke 'print' in production code. Try using a logging framework -
          lib\widgets\payment_dialog.dart:119:7 - avoid_print
warning - This class (or a class that this class inherits from) is marked as '@immutable', but one or
       more of its instance fields aren't final: PredictionPlaceUI.predictedPlaceData -
       lib\widgets\prediction_place_ui.dart:10:7 - must_be_immutable
   info - Missing type annotation. Try adding a type annotation -
          lib\widgets\prediction_place_ui.dart:20:3 - strict_top_level_inference
   info - Don't invoke 'print' in production code. Try using a logging framework -
          lib\widgets\prediction_place_ui.dart:32:5 - avoid_print
   info - Don't use 'BuildContext's across async gaps. Try rewriting the code to not use the
          'BuildContext', or guard the use with a 'mounted' check -
          lib\widgets\prediction_place_ui.dart:39:19 - use_build_context_synchronously
   info - Don't invoke 'print' in production code. Try using a logging framework -
          lib\widgets\prediction_place_ui.dart:43:7 - avoid_print
   info - Don't invoke 'print' in production code. Try using a logging framework -
          lib\widgets\prediction_place_ui.dart:48:5 - avoid_print
   info - Don't invoke 'print' in production code. Try using a logging framework -
          lib\widgets\prediction_place_ui.dart:56:7 - avoid_print
   info - Don't invoke 'print' in production code. Try using a logging framework -
          lib\widgets\prediction_place_ui.dart:62:7 - avoid_print
   info - Don't invoke 'print' in production code. Try using a logging framework -
          lib\widgets\prediction_place_ui.dart:68:7 - avoid_print
   info - Don't invoke 'print' in production code. Try using a logging framework -
          lib\widgets\prediction_place_ui.dart:73:7 - avoid_print
   info - Don't use 'BuildContext's across async gaps. Try rewriting the code to not use the
          'BuildContext', or guard the use with a 'mounted' check -
          lib\widgets\prediction_place_ui.dart:76:33 - use_build_context_synchronously
   info - Don't use 'BuildContext's across async gaps. Try rewriting the code to not use the
          'BuildContext', or guard the use with a 'mounted' check -
          lib\widgets\prediction_place_ui.dart:80:21 - use_build_context_synchronously
   info - Don't invoke 'print' in production code. Try using a logging framework -
          lib\widgets\prediction_place_ui.dart:82:7 - avoid_print
   info - Constructors in '@immutable' classes should be declared as 'const'. Try adding 'const' to
          the constructor declaration - lib\widgets\sign_out_dialog.dart:8:3 -
          prefer_const_constructors_in_immutables

153 issues found. (ran in 11.3s)