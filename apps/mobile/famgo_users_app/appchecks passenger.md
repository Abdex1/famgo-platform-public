https://developer.apple.com/documentation/devicecheck/establishing-your-app-s-integrity
Step 2: Register iOS App (com.famgo.famgoPassenger)

For iOS, you can use App Attest (recommended for iOS 14+) or DeviceCheck.

Apple Developer Portal: Enable App Attest in your App ID configuration. For DeviceCheck, generate and download a .p8 private key, and note your Key ID and Team ID.

Register Provider: In the Firebase Console > App Check > Apps, select your iOS app:

App Attest: Add your Apple Team ID.

DeviceCheck: Upload the .p8 key and input your Key ID and Team ID.

Step 3: Register Web Apps

Secure your web clients using reCAPTCHA Enterprise or reCAPTCHA v3.

Create Keys: Register your authorized domains (e.g., localhost, famgo-21a0d.web.app, and famgo-21a0d.firebaseapp.com) on the Google Cloud or reCAPTCHA Admin console to get a Site Key and a Secret Key.

Add to Firebase Console: In App Check > Apps, select your Web app, choose reCAPTCHA, and input the Secret Key.

CLI/Environment Configuration: Save your generated site key in your project's .env environment variables [7]:

NEXT_PUBLIC_RECAPTCHA_SITE_KEY=your_recaptcha_site_key
Because your famgo project is on the Spark billing plan, registering these providers operates under a no-cost tier, though some external APIs (like reCAPTCHA Enterprise) have their own limits.
##########################################################################################
https://developer.apple.com/help/account/capabilities/create-a-devicecheck-private-key/

https://developer.apple.com/help/account/keys/create-a-private-key
Create a DeviceCheck private key
To authenticate communication with the DeviceCheck service, you’ll use a private key enabled with DeviceCheck.

First create and download a private key with DeviceCheck enabled. Then get the key identifier (kid) to create a JSON Web Token (JWT) that you’ll use to communicate with the capabilities you enabled. To use this identifier in a DeviceCheck JWT, visit Accessing and Modifying Per-Device Data.

If you suspect a private key is compromised, first create a new private key with DeviceCheck enabled. Then, after transitioning to the new key, revoke the old private key.

To learn about the DeviceCheck API, visit DeviceCheck.

Required role: Account Holder or Admin.

####################################################################

https://www.google.com/recaptcha/admin/site/757145669/setup
Adding reCAPTCHA to your site
'famgo-21a0d.passenger.web.app' has been registered.
Use this site key in the HTML code your site serves to users.
6LdFICEtAAAAANjRlMuPrDYgc6a1bJ1s6HANIzo7

Server side integration

https://developers.google.com/recaptcha/docs/invisible
Use this secret key for communication between your site and reCAPTCHA.
6LdFICEtAAAAAHyfDMJ8vbhtcgZRUmItYWx0gZL9


Server side integration

https://developers.google.com/recaptcha/docs/invisible
###############################################################################



site key Server side integration


Invisible reCAPTCHA

Spark icon
Page Summary
This page explains how to enable and customize the invisible reCAPTCHA on your webpage.

To invoke the invisible reCAPTCHA, you can either:

Automatically bind the challenge to a button or
Programmatically bind the challenge to a button or
Programmatically invoke the challenge
See Configurations to learn how to customize the invisible reCAPTCHA. For example, you may want to specify the language or badge location.

See Verifying the user's response to check if the user successfully solved the CAPTCHA.

Automatically bind the challenge to a button
The easiest method for using the invisible reCAPTCHA widget on your page is to include the necessary JavaScript resource and add a few attributes to your html button. The necessary attributes are a class name 'g-recaptcha', your site key in the data-sitekey attribute, and the name of a JavaScript callback to handle completion of the captcha in the data-callback attribute.


<html>
  <head>
    <title>reCAPTCHA demo: Simple page</title>
     <script src="https://www.google.com/recaptcha/api.js" async defer></script>
     <script>
       function onSubmit(token) {
         document.getElementById("demo-form").submit();
       }
     </script>
  </head>
  <body>
    <form id="demo-form" action="?" method="POST">
      <button class="g-recaptcha" data-sitekey="your_site_key" data-callback="onSubmit">Submit</button>
      <br/>
    </form>
  </body>
</html>
The script must be loaded using the HTTPS protocol and can be included from any point on the page without restriction.

Programmatically bind the challenge to a button or invoke the challenge.
Deferring the binding can be achieved by specifying your onload callback function and adding parameters to the JavaScript resource. This works the same as the normal reCAPTCHA challenge.

Programmatically invoke the challenge.
Invoking the reCAPTCHA verification programmatically can be achieved by rendering the challenge in a div with an attribute data-size="invisible" and programmatically calling execute.

Create a div with data-size="invisible".


<div class="g-recaptcha"
      data-sitekey="_your_site_key_"
      data-callback="onSubmit"
      data-size="invisible">
</div>
Call grecaptcha.execute from a javascript method.


grecaptcha.execute();
When your callback is executed, you can call the grecaptcha.render method from the JavaScript API.

Your onload callback function must be defined before the reCAPTCHA API loads. To ensure there are no race conditions:
order your scripts with the callback first, and then reCAPTCHA
use the `async` and `defer` parameters in the `script` tags
Configuration
JavaScript resource (api.js) parameters
Parameter	Value	Description
onload		Optional. The name of your callback function to be executed once all the dependencies have loaded.
render	explicit
onload	Optional. Whether to render the widget explicitly. Defaults to onload, which will render the widget in the first g-recaptcha tag it finds.
hl	See language codes	Optional. Forces the widget to render in a specific language. Auto-detects the user's language if unspecified.
g-recaptcha tag attributes and grecaptcha.render parameters
g-recaptcha tag attribute	grecaptcha.render parameter	Value	Default	Description
data-sitekey	sitekey			Your sitekey.
data-badge	badge	bottomright bottomleft inline	bottomright	Optional. Reposition the reCAPTCHA badge. 'inline' lets you position it with CSS.
data-size	size	invisible		Optional. Used to create an invisible widget bound to a div and programmatically executed.
data-tabindex	tabindex		0	Optional. The tabindex of the challenge. If other elements in your page use tabindex, it should be set to make user navigation easier.
data-callback	callback			Optional. The name of your callback function, executed when the user submits a successful response. The g-recaptcha-response token is passed to your callback.
data-expired-callback	expired-callback			Optional. The name of your callback function, executed when the reCAPTCHA response expires and the user needs to re-verify.
data-error-callback	error-callback			Optional. The name of your callback function, executed when reCAPTCHA encounters an error (usually network connectivity) and cannot continue until connectivity is restored. If you specify a function here, you are responsible for informing the user that they should retry.
isolated		false	Optional. For plugin owners to not interfere with existing reCAPTCHA installations on a page. If true, this reCAPTCHA instance will be part of a separate ID space.
JavaScript API
Method	Description
grecaptcha.render (
container,
parameters,
inherit
)	Renders the container as a reCAPTCHA widget and returns the ID of the newly created widget.
container
  The HTML element to render the reCAPTCHA widget.  Specify either the ID of the container (string) or the DOM element itself.
parameters
  An object containing parameters as key=value pairs, for example, {"sitekey": "your_site_key", "theme": "light"}. See grecaptcha.render parameters.
inherit
  Use existing data-* attributes on the element if the corresponding parameter is not specified. The parameters will take precedence over the attributes.
grecaptcha.execute(
opt_widget_id
)	Programmatically invoke the reCAPTCHA check. Used if the invisible reCAPTCHA is on a div instead of a button.
opt_widget_id
  Optional widget ID, defaults to the first widget created if unspecified.
grecaptcha.reset(
opt_widget_id
)	Resets the reCAPTCHA widget.
opt_widget_id
  Optional widget ID, defaults to the first widget created if unspecified.
grecaptcha.getResponse(
opt_widget_id
)	Gets the response for the reCAPTCHA widget.
opt_widget_id
  Optional widget ID, defaults to the first widget created if unspecified.
Examples
Explicit rendering after an onload callback


<html>
  <head>
    <title>reCAPTCHA demo: Explicit render after an onload callback</title>
    <script>
        var onSubmit = function(token) {
          console.log('success!');
        };

        var onloadCallback = function() {
          grecaptcha.render('submit', {
            'sitekey' : 'your_site_key',
            'callback' : onSubmit
          });
        };
    </script>
  </head>
  <body>
    <form action="?" method="POST">
      <input id="submit" type="submit" value="Submit">
    </form>
    <script src="https://www.google.com/recaptcha/api.js?onload=onloadCallback&render=explicit"
        async defer>
    </script>
  </body>
</html>
Invoking the invisible reCAPTCHA challenge after client side validation.


<html>
  <head>
  <script>
    function onSubmit(token) {
      alert('thanks ' + document.getElementById('field').value);
    }

    function validate(event) {
      event.preventDefault();
      if (!document.getElementById('field').value) {
        alert("You must add text to the required field");
      } else {
        grecaptcha.execute();
      }
    }

    function onload() {
      var element = document.getElementById('submit');
      element.onclick = validate;
    }
  </script>
  <script src="https://www.google.com/recaptcha/api.js" async defer></script>
  </head>
  <body>
    <form>
      Name: (required) <input id="field" name="field">
      <div id="recaptcha" class="g-recaptcha"
          data-sitekey="_your_site_key_"
          data-callback="onSubmit"
          data-size="invisible"></div>
      <button id="submit">submit</button>
    </form>
    <script>onload();</script>
  </body>
</html>
##########################################################

Secret key Server side integration for communication between your site and reCAPTCHA

Verifying the user's response



Spark icon
Page Summary
This page explains how to verify a user's response to a reCAPTCHA challenge from your application's backend.

For web users, you can get the user’s response token in one of three ways:

g-recaptcha-response POST parameter when the user submits the form on your site
grecaptcha.getResponse(opt_widget_id) after the user completes the reCAPTCHA challenge
As a string argument to your callback function if data-callback is specified in either the g-recaptcha tag attribute or the callback parameter in the grecaptcha.render method
For Android library users, you can call the SafetyNetApi.RecaptchaTokenResult.getTokenResult() method to get response token if the status returns successful.

Token Restrictions
Each reCAPTCHA user response token is valid for two minutes, and can only be verified once to prevent replay attacks. If you need a new token, you can re-run the reCAPTCHA verification.

After you get the response token, you need to verify it within two minutes with reCAPTCHA using the following API to ensure the token is valid.

API Request
URL: https://www.google.com/recaptcha/api/siteverify

METHOD: POST

POST Parameter	Description
secret	Required. The shared key between your site and reCAPTCHA.
response	Required. The user response token provided by the reCAPTCHA client-side integration on your site.
remoteip	Optional. The user's IP address.
API Response
The response is a JSON object:


{
  "success": true|false,
  "challenge_ts": timestamp,  // timestamp of the challenge load (ISO format yyyy-MM-dd'T'HH:mm:ssZZ)
  "hostname": string,         // the hostname of the site where the reCAPTCHA was solved
  "error-codes": [...]        // optional
}
For reCAPTCHA Android:


{
  "success": true|false,
  "challenge_ts": timestamp,  // timestamp of the challenge load (ISO format yyyy-MM-dd'T'HH:mm:ssZZ)
  "apk_package_name": string, // the package name of the app where the reCAPTCHA was solved
  "error-codes": [...]        // optional
}
Error code reference
Error code	Description
missing-input-secret	The secret parameter is missing.
invalid-input-secret	The secret parameter is invalid or malformed.
missing-input-response	The response parameter is missing.
invalid-input-response	The response parameter is invalid or malformed.
bad-request	The request is invalid or malformed.
timeout-or-duplicate	The response is no longer valid: either is too old or has been used previously.