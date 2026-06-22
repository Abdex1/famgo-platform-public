cd android
.\gradlew assembleDebug --stacktrace



If Flutter still references old plugins, usually one of these happened:

Android build cache is corrupted.
.flutter-plugins-dependencies is stale.
Old generated files remain in Android.
Old embedding code exists somewhere.
First check

Run:

Get-ChildItem -Recurse -Include *.java,*.kt |
Select-String "ImageCropperPlugin"

If you find ONLY:

android/app/src/main/java/io/flutter/plugins/GeneratedPluginRegistrant.java

then it is a stale generated file.

Delete all generated Android files

From project root:

Remove-Item .dart_tool -Recurse -Force
Remove-Item build -Recurse -Force
Remove-Item android\.gradle -Recurse -Force

Remove-Item .flutter-plugins -Force -ErrorAction SilentlyContinue
Remove-Item .flutter-plugins-dependencies -Force -ErrorAction SilentlyContinue
Delete GeneratedPluginRegistrant

Delete:

android/app/src/main/java/io/flutter/plugins/

Entire folder:

Remove-Item `
"android\app\src\main\java\io\flutter\plugins" `
-Recurse `
-Force


If you want to replace every occurrence inside code (imports, variable names, class names, strings, comments, etc.) and not just filenames, use a content replacement script.

First, preview how many matches exist
Get-ChildItem "C:\Users\FEMOS\Desktop\Femos\extrac\famGo\famgo_driver_app\lib" -Recurse -Filter "*.dart" |
Select-String -Pattern "uber_drivers_app"
Then replace everywhere inside .dart files
Get-ChildItem "C:\Users\FEMOS\Desktop\Femos\extrac\famGo\famgo_driver_app\lib" -Recurse -Filter "*.dart" | ForEach-Object {
    $content = Get-Content $_.FullName -Raw
    $content = $content.Replace("uber_drivers_app", "famgo_drivers_app")
    Set-Content $_.FullName $content -Encoding UTF8
}

This changes things like:

import 'package:uber_drivers_app/global/global.dart';

↓

import 'package:famgo_drivers_app/global/global.dart';

and

String appName = "uber_drivers_app";

↓

String appName = "famgo_drivers_app";
If you also want to replace PascalCase names

For example:

UberDriversApp

↓

FamgoDriversApp

Run:

Get-ChildItem "C:\Users\FEMOS\Desktop\Femos\extrac\famGo\famgo_driver_app\lib" -Recurse -Filter "*.dart" | ForEach-Object {
    $content = Get-Content $_.FullName -Raw
    $content = $content.Replace("UberDriversApp", "FamgoDriversApp")
    Set-Content $_.FullName $content -Encoding UTF8
}
If you want to replace all common variants in one go
Get-ChildItem "C:\Users\FEMOS\Desktop\Femos\extrac\famGo\famgo_driver_app\lib" -Recurse -Filter "*.dart" | ForEach-Object {

    $content = Get-Content $_.FullName -Raw

    $content = $content.Replace("uber_drivers_app", "famgo_drivers_app")
    $content = $content.Replace("UberDriversApp", "FamgoDriversApp")
    $content = $content.Replace("uberDriversApp", "famgoDriversApp")
    $content = $content.Replace("UBER_DRIVERS_APP", "FAMGO_DRIVERS_APP")

    Set-Content $_.FullName $content -Encoding UTF8
}
Verify nothing remains
Get-ChildItem "C:\Users\FEMOS\Desktop\Femos\extrac\famGo\famgo_driver_app\lib" -Recurse -Filter "*.dart" |
Select-String -Pattern "uber_drivers_app|UberDriversApp|uberDriversApp|UBER_DRIVERS_APP"

If no results are returned, all matching identifiers, imports, variables, strings, and comments inside lib have been replaced.