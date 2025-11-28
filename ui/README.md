# UI Component: RedactionDashboard

This directory contains the source code for the `RedactionDashboard` React component, which serves as the main user interface for the Cilium-Shield project.

## Overview

The `RedactionDashboard` provides a "CISO Command Center" view of all redaction events detected by the Cilium-Shield system. It fetches data from the `shield-observer` control plane and displays it in a user-friendly and informative way.

## Features

*   **Live Redaction Feed:** A real-time table that displays the latest redaction events, including the timestamp, source pod IP, the type of data that was redacted, and the destination of the request.
*   **Total Redactions:** A prominent display of the total number of redactions that have occurred.
*   **Redactions by Type:** A breakdown of the number of redactions for each type of sensitive data (e.g., credit card numbers, API keys, email addresses).
*   **Dark Theme:** A visually appealing dark theme that is suitable for a security-focused dashboard.
*   **Auto-Refreshing:** The dashboard automatically fetches new data every 5 seconds to provide a near real-time view of the system's activity.
*   **Error Handling:** The dashboard displays a clear error message if it is unable to connect to the `shield-observer` control plane.

## Getting Started

To run the `RedactionDashboard` component, you will need to have a React development environment set up.

1.  **Install Dependencies:**
    ```sh
    npm install
    ```
2.  **Run the Development Server:**
    ```sh
    npm run dev
    ```
    This will start the development server and open the dashboard in your browser.

## Component Breakdown

*   **`RedactionDashboard.jsx`:** The main React component that renders the dashboard.
    *   **State Management:** The component uses the `useState` and `useEffect` hooks to manage its state, including the list of events, loading status, and error messages.
    *   **Data Fetching:** The `fetchEvents` function is responsible for fetching data from the `shield-observer` control plane.
    *   **UI Rendering:** The component uses Tailwind CSS for styling and renders the different sections of the dashboard, including the header, summary statistics, and the live event feed.
