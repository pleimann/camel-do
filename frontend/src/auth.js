// Google One Tap Authentication Handler

/**
 * Initializes Google One Tap authentication
 * @param {string} clientId - Google OAuth client ID
 * @param {Function} onSuccess - Callback function when authentication succeeds
 * @param {Function} onError - Callback function when authentication fails
 */
export function initializeGoogleOneTap(clientId, onSuccess, onError) {
  // Check if Google One Tap is already initialized
  if (window.google && document.getElementById('g_id_onload')) {
    console.log('Google One Tap already initialized');
    return;
  }

  // Create the Google One Tap container
  const container = document.createElement('div');
  container.id = 'g_id_onload';
  container.setAttribute('data-client_id', clientId);
  container.setAttribute('data-context', 'signin');
  container.setAttribute('data-ux_mode', 'popup');
  container.setAttribute('data-callback', 'handleGoogleOneTapResponse');
  container.setAttribute('data-auto_prompt', 'true');
  document.body.appendChild(container);

  // Create the Google Sign In button container
  const buttonContainer = document.createElement('div');
  buttonContainer.className = 'g_id_signin';
  buttonContainer.setAttribute('data-type', 'standard');
  buttonContainer.setAttribute('data-shape', 'rectangular');
  buttonContainer.setAttribute('data-theme', 'outline');
  buttonContainer.setAttribute('data-text', 'signin_with');
  buttonContainer.setAttribute('data-size', 'large');
  buttonContainer.setAttribute('data-logo_alignment', 'left');
  
  // Add the button to a hidden div
  const hiddenDiv = document.createElement('div');
  hiddenDiv.style.display = 'none';
  hiddenDiv.appendChild(buttonContainer);
  document.body.appendChild(hiddenDiv);

  // Define the callback function
  window.handleGoogleOneTapResponse = (response) => {
    if (response && response.credential) {
      // Send the credential to the backend
      fetch('/auth/google/onetap', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: `credential=${encodeURIComponent(response.credential)}`,
      })
        .then(res => res.json())
        .then(data => {
          if (data.status === 'success') {
            if (onSuccess) onSuccess(response);
          } else {
            if (onError) onError(new Error('Authentication failed'));
          }
        })
        .catch(err => {
          console.error('Error during authentication:', err);
          if (onError) onError(err);
        });
    } else {
      console.error('Invalid response from Google One Tap');
      if (onError) onError(new Error('Invalid response from Google One Tap'));
    }
  };

  // Load the Google Identity Services script if not already loaded
  if (!window.google) {
    const script = document.createElement('script');
    script.src = 'https://accounts.google.com/gsi/client';
    script.async = true;
    script.defer = true;
    document.head.appendChild(script);
  }
}

/**
 * Signs out the user from Google One Tap
 */
export function signOut() {
  if (window.google) {
    google.accounts.id.disableAutoSelect();
    // Remove any stored credentials
    document.cookie = 'g_csrf_token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
    localStorage.removeItem('googleOneTapUser');
  }
}

/**
 * Checks if the user is authenticated with Google One Tap
 * @returns {boolean} True if authenticated, false otherwise
 */
export function isAuthenticated() {
  return !!localStorage.getItem('googleOneTapUser');
}

/**
 * Gets the authenticated user from local storage
 * @returns {Object|null} User object or null if not authenticated
 */
export function getAuthenticatedUser() {
  const userJson = localStorage.getItem('googleOneTapUser');
  return userJson ? JSON.parse(userJson) : null;
}
