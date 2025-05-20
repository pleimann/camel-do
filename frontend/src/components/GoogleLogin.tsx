import { Component, createSignal } from "solid-js";

const GoogleSignIn: Component = () => {
  const [isSignedIn, setIsSignedIn] = createSignal(false);
  const [user, setUser] = createSignal<any>(null);

  // Function to handle Google Sign-In
  const handleGoogleSignIn = async () => {
    // Call the Google Sign-In API
    const response = await google.accounts.id.login(
      {
        client_id: "YOUR_CLIENT_ID", // Replace with your client ID
        callback: (resp) => {
          // Handle the response from Google Sign-In
          if (resp.credential) {
            // Access the user's profile information
            // Example: Get the user's ID token
            const idToken = resp.credential;
            console.log("ID Token:", idToken);

            // You can use the ID token to fetch user data from your backend
            // or use a library like Firebase to handle authentication

            // Update the state to reflect that the user is signed in
            setIsSignedIn(true);
            // Fetch user data
            fetchUser(idToken);
          } else {
            console.error("Sign-in failed:", resp);
          }
        },
      }
    );
  };

  const fetchUser = async (idToken: string) => {
    try {
      const response = await fetch('/api/user', { // Replace with your backend endpoint
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ idToken: idToken })
      });
      const data = await response.json();
      setUser(data);
    } catch (error) {
      console.error("Error fetching user data:", error);
    }
  };

  // Function to handle Google Sign-Out
  const handleGoogleSignOut = () => {
    // Call the Google Sign-Out API
    google.accounts.id.signOut();
    // Update the state to reflect that the user is signed out
    setIsSignedIn(false);
    setUser(null);
  };

  return (
    <>
      {isSignedIn() ? (
        <>
          <p>Signed in as: {user()?.email}</p>
          <button onClick={handleGoogleSignOut}>Sign Out</button>
        </>
      ) : (
        <button onClick={handleGoogleSignIn}>Sign In with Google</button>
      )}
    </>
  );
};

export default GoogleSignIn;
