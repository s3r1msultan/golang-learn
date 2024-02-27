# Restaurant Online System

## Participants

- Adilkhan Yeslyambek
- Yersultan Serimbetov
- Tayir Taishanov

## Overview

This project is an online restaurant system designed to provide a seamless user experience in ordering food, managing user profiles, accessing order history, earning bonuses, setting delivery addresses, and contacting customer support. It aims to offer a user-friendly interface for enhanced engagement and satisfaction.

### Technologies Used

- **HTML & CSS**
- **JavaScript**
- **Bootstrap**
- **jQuery**

## Table of Contents

1. [Main Page (index.html)](#main-page)
2. [Sign In Page (SignInPage)](#sign-in-page)
3. [User Profile Page (UserPage)](#user-profile-page)
4. [Order History Page (OrderHistoryPage)](#order-history-page)
5. [Bonus Page (BonusPage)](#bonus-page)
6. [Delivery Address Page (DeliveryAddressPage)](#delivery-address-page)
7. [Contact Support Page (SupportPage)](#contact-support-page)
8. [Server-side (main.go)](#server-side)

## Pages Description

### Main Page (index.html)

- **Navigation:** Bootstrap navbar with links to "Menu" and "Sign In" pages.
- **Design Elements:** Utilizes Bootstrap cards for displaying food items.
- **Carousel:** Incorporates a Bootstrap carousel displaying restaurant ambiance.
- **Multimedia:** Embeds two videos, one using JavaScript to enable fullscreen mode.

### Sign In Page (SignInPage)

- **Authentication Form:** Contains input fields for username and password.
- **JavaScript Functionality:** Validates user input before submitting the form.
- **Redirects:** Upon successful authentication, directs authenticated users to the "User Profile" page.
- **Bootstrap Styling:** Utilizes Bootstrap for form layout and styling.

### User Profile Page (UserPage)

- **Profile Display:** Shows user details like name, email, and preferences.
- **Editable Fields:** Allows users to update their profile information.
- **Form Elements:** Utilizes Bootstrap's form components for an organized layout and provides an interactive To-Do list.

### Order History Page (OrderHistoryPage)

- **Tabular Data:** Displays previous orders in an HTML table format.
- **Order Details:** Includes date, ordered dishes, and total price columns.
- **Styling:** Applies Bootstrap styles for the table, enhancing readability.
- **Interactive Links:** Provides a "Reorder" link for each order.

### Bonus Page (BonusPage)

- **Rewards Listing:** Lists earned bonuses through rewards programs.
- **Table Layout:** Uses Bootstrap table components for a structured display.
- **Information Display:** Shows details like date, bonus descriptions, and points earned.

### Delivery Address Page (DeliveryAddressPage)

- **Form Input:** Includes fields for full name, address, city, zip code, and phone number.
- **Form Submission:** Allows users to save their delivery addresses.
- **Bootstrap Form Styling:** Utilizes Bootstrap's form elements and styles for a polished look.
- **Validation:** Ensures required fields are filled out before submission.

### Contact Support Page (SupportPage)

- **Support Form:** Provides fields for full name, email, and message.
- **Validation Logic:** Implements JavaScript validation for form inputs.
- **Loading Spinner:** Displays a loading spinner while submitting the form.
- **Modals:** Uses Bootstrap modals for alerting successful form submissions and chat functionality.

## JavaScript and Bootstrap Usage

- **JavaScript (JS):**
  - Handles form submissions, validations, and redirects.
  - Controls dynamic content updates and chat functionality.
  - Manages modals, alerts, and loading spinners.
  
- **Bootstrap:**
  - Ensures responsive design across different devices.
  - Provides pre-styled components for UI elements.
  - Enhances visual appeal and user experience.

## Server-side (main.go)

### Overview

The server-side of the Restaurant Online System is implemented in Go (Golang) using the main.go file. This server handles HTTP requests, interacts with the MongoDB database, and provides the necessary endpoints for the client-side functionalities.

### Technologies Used

- **Go (Golang):**
  - Used for building the server-side logic.
  - Enables efficient handling of HTTP requests and responses.
  - Integrates with the MongoDB driver for database interactions.

- **MongoDB:**
  - Utilized as the database for storing and retrieving delivery address information.
  - Integrated with the Go application using the MongoDB Go driver.

### Main Functionality

- **HTTP Server:**
  - Creates an HTTP server to listen for incoming requests.
  - Defines endpoints for handling delivery address data.

- **MongoDB Integration:**
  - Connects to the MongoDB database to store and retrieve delivery address information.

- **Endpoints:**
  - `/` - Handles requests for the main page.
  - `/delivery` - Handles requests related to delivery addresses, including creation and retrieval.

- **Error Handling:**
  - Implements error handling for various scenarios, such as invalid JSON format and MongoDB connection issues.

### Usage

- **Insert Data to MongoDB:**
  - Handles POST requests to `/delivery` for creating new delivery addresses.
  - Validates incoming JSON data and inserts it into the MongoDB collection.

- **Retrieve Data from MongoDB:**
  - Handles GET requests to `/delivery` for retrieving a list of all delivery addresses.

### Note

- Ensure that the MongoDB connection URI is correctly configured in the `insertDataToMongoDB` function.
- The provided `main.go` file is a basic example and might need additional improvements for production use, such as error logging and security considerations.

