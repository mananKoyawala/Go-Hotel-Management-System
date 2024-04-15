# Hotel Management System 🏨

Hotel management system is Backend project completely written in **GO (Golang)**. It aims to provide a robust and efficient solution for managing various aspects of hotel operations.

## Database and Cloud Services Integration

1. **MongoDB**:

   - Used MongoDB as the database for the hotel management system due to its NoSQL nature and ease of building APIs.
   - Learn more about how to integrate MongoDB [here](https://docs.mongodb.com/).

2. **Firebase**:
   - Utilized Firebase as an image storage solution for storing guest profiles, room photos, and other image assets.
   - Explore Firebase features and integrations [here](https://firebase.google.com/).

## Authentication and Communication Services

1. **JWT Authentication**:

   - Implemented JWT authentication with four access types: Admin, Manager, Guest, and Driver. Each access type has specific permissions and privileges within the system.
   - In this system I implemented jwt token as token(24 hours valid from login) and refersh_token(168 hours valid from login).
   - Learn more about how JWT [here](https://jwt.io/introduction).

2. **Resend**:
   - Integrated an email resend service for sending email verification links to guests during the registration process. This enhances user experience and ensures email verification completion.
   - Explore Resend email service [here](http://resend.com).

## Roles

- **Admin** :- The Admin oversees all operations within the hotel management system and has access to all functionalities.
- **Manager** :- Each branch of the hotel is managed by a Manager who is responsible for overseeing operations within their branch.
- **Driver** :- Drivers provide pickup services for guests, ensuring smooth transportation arrangements during their stay.
- **Guest** :- Guests are the primary users of the system, utilizing its services for booking, managing reservations, and accessing various amenities offered by the hotel.

## Installation

1. Ensure you have Go installed on your machine. If not, you can download it [here](https://golang.org/dl/).

2. Clone this repository to your local machine:

```bash
   git clone https://github.com/mananKoyawala/Go-Hotel-Management-System.git
```

3. Navigate to the project directory:

```bash
   cd Go-Hotel-Management-System
```

4. Install required packages: Install the necessary external packages using the following commands.

- For GIN framework :

```bash
   go get github.com/gin-gonic/gin
```

- For use color package:

```bash
   go get github.com/fatih/color
```

- For handling environment variables:

```bash
   go get github.com/joho/godotenv
```

- For JWT token-based authentication:

```bash
   go get github.com/dgrijalva/jwt-go
```

- For using the Resend email service:

```bash
   go get github.com/resend/resend-go/v2
```

- For establishing connection with Firebase and Firebase storage:

```bash
   go get cloud.google.com/go/storage
   go get firebase.google.com/go
```

5. Make additional changes to the code (essentiall):

   - Replace the serviceAccountKey.json file in the _pkg/service/image-upload/Image-upload-helper.go_ folder with your own Firebase project service account key.
   - Edit the .env file with your own details.
   - Change you domain for sending the verification emails in _pkg/service/Email-Verification/Email-Verification-Service.go_

6. Run the game using the following command:

```bash
   go run main.go
```

## API Documentation

- **Note**: The following routes are intended for use on a localhost server. Access tokens are required for protected routes. Please include the access token as a header named 'X-Auth-Token' in your requests. Tokens are obtained upon user login.

### Admin Route

- **Admin Login**:
  - **Method**: POST
  - **Endpoint**: `http://localhost:8000/admin/login`
  - **Access**: Public
  - **Parameters**: email and password (Form Data)

### Branch Routes

- **Get all Branches**:

  - **Method**: GET
  - **Endpoint**: `http://localhost:8000/branch/getall`
  - **Access**: Public
  - **Pagination**: `recordPerPage` (default: 10), `page` (default: 1)

- **Get one Branch**:

  - **Method**: GET
  - **Endpoint**: `http://localhost:8000/branch/get/:id`
  - **Access**: Public
  - **Parameters**: branch_id

- **Get all Branches by Status**:

  - **Method**: GET
  - **Endpoint**: `http://localhost:8000/branch/get-branch-by-status/:status`
  - **Access**: Admin Only
  - **Parameters**: status (0 or 1)
  - **Pagination**: `recordPerPage` (default: 10), `page` (default: 1)

- **Create Branch**:

  - **Method**: POST
  - **Endpoint**: `http://localhost:8000/branch/create`
  - **Access**: Admin Only
  - **Data**: Form Data
  - **Parameters**: manager_id, branch_name, address, phone, email, city, state, country, pincode, file

- **Update Branch Details**:

  - **Method**: PUT
  - **Endpoint**: `http://localhost:8000/branch/update-all/:id`
  - **Access**: Admin Only
  - **Data**: Form Data
  - **Parameters**: manager_id, branch_name, address, phone, email, city, state, country, pincode

- **Update Branch Status**:

  - **Method**: PATCH
  - **Endpoint**: `http://localhost:8000/branch/update-branch-status/:id`
  - **Access**: Admin Only
  - **Note**: It toggles the status

- **Add Image to Branch**:

  - **Method**: PATCH
  - **Endpoint**: `http://localhost:8000/branch/add-image/:id`
  - **Access**: Admin Only
  - **Data**: Form Data
  - **Parameters**: file

- **Delete Image from Branch**:

  - **Method**: DELETE
  - **Endpoint**: `http://localhost:8000/branch/delete-image/:id`
  - **Access**: Admin Only
  - **Data**: Form Data
  - **Parameters**: image (uploaded image url)

- **Delete Branch**:

  - **Method**: DELETE
  - **Endpoint**: `http://localhost:8000/branch/delete/:id`
  - **Access**: Admin Only

- **Search Branch**:

  - **Method**: POST
  - **Endpoint**: `http://localhost:8000/branch/search`
  - **Access**: Public
  - **Search by**: Branch Name, Address, Status, City, State, Country as search (Form Field)
  - **Pagination**: `recordPerPage` (default: 10), `page` (default: 1)

- **Filter Branch**:

  - **Method**: POST
  - **Endpoint**: `http://localhost:8000/branch/filter`
  - **Access**: Public
  - **Filter by**: city, state, country, status (as Form Field)
  - **Pagination**: `recordPerPage` (default: 10), `page` (default: 1)

### Room Routes

- **Get all Rooms**:

  - **Method**: GET
  - **Endpoint**: `http://localhost:8000/room/getall`
  - **Access**: Public
  - **Pagination**: `recordPerPage` (default: 10), `page` (default: 1)

- **Get all Rooms by Branch**:

  - **Method**: GET
  - **Endpoint**: `http://localhost:8000/room/getall/:id`
  - **Access**: Public
  - **Pagination**: `recordPerPage` (default: 10), `page` (default: 1)

- **Get one Room**:

  - **Method**: GET
  - **Endpoint**: `http://localhost:8000/room/get/:id`
  - **Access**: Public

- **Create Room**:

  - **Method**: POST
  - **Endpoint**: `http://localhost:8000/room/create`
  - **Access**: Manager Only
  - **Data**: Form Data
  - **Parameters**: branch_id, room_number, room_type, price, capacity, file

- **Update Room**:

  - **Method**: PUT
  - **Endpoint**: `http://localhost:8000/room/update-all/:id`
  - **Access**: Manager Only
  - **Data**: Form Data
  - **Parameters**: room_number, room_type, cleaning_status, room_availability, price

- **Add Image to Room**:

  - **Method**: PATCH
  - **Endpoint**: `http://localhost:8000/room/add-image/:id`
  - **Access**: Manager Only
  - **Data**: Form Data
  - **Parameters**: file

- **Delete Image from Room**:

  - **Method**: DELETE
  - **Endpoint**: `http://localhost:8000/room/delete-image/:id`
  - **Access**: Manager Only
  - **Data**: Form Data
  - **Parameters**: image (uploaded image url)

- **Delete Room**:

  - **Method**: DELETE
  - **Endpoint**: `http://localhost:8000/room/delete/:id`
  - **Access**: Manager Only

- **Filter Room**:

  - **Method**: POST
  - **Endpoint**: `http://localhost:8000/room/filter`
  - **Access**: Public
  - **Data** : Form Data
  - **Filter by**: room_type, room_availability, cleaning_status, price (as Form Field)
  - **Pagination**: `recordPerPage` (default: 10), `page` (default: 1)

### Manager Routes

- **Manager Login**:

  - **Method**: POST
  - **Endpoint**: `http://localhost:8000/manager/login`
  - **Access**: Public
  - **Data** : Form Data
  - **Paramaters** : email and password

- **Get all Managers**:

  - **Method**: GET
  - **Endpoint**: `http://localhost:8000/manager/getall`
  - **Access**: Admin Only
  - **Pagination**: `recordPerPage` (default: 10), `page` (default: 1)

- **Get Manager**:

  - **Method**: GET
  - **Endpoint**: `http://localhost:8000/manager/get/:id`
  - **Access**: Admin Only

- **Create Manager**:

  - **Method**: POST
  - **Endpoint**: `http://localhost:8000/manager/create`
  - **Access**: Admin Only
  - **Data** : Form Data
  - **Paramaters** : first_name, last_name, age, phone, email, password, gender, salary, aadhar_number, file

- **Update Manager Status**:

  - **Method**: PATCH
  - **Endpoint**: `http://localhost:8000/manager/update-status/:id`
  - **Access**: Admin Only
  - **Note**: It toggles the status

- **Update Manager Details**:

  - **Method**: PUT
  - **Endpoint**: `http://localhost:8000/manager/update-all/:id`
  - **Access**: Admin Only
  - **Data** : Form Data
  - **Paramaters** : first_name, last_name, age, gender, salary, phone

- **Delete Manager**:

  - **Method**: DELETE
  - **Endpoint**: `http://localhost:8000/manager/delete/:id`
  - **Access**: Admin Only

- **Reset Manager Password**:

  - **Method**: PATCH
  - **Endpoint**: `http://localhost:8000/manager/update-password`
  - **Access**: Admin Only
  - **Data** : Form Data
  - **Paramaters** : email and password

- **Update Manager Profile Picture**:

  - **Method**: PATCH
  - **Endpoint**: `http://localhost:8000/manager/update-profile-pic/:id`
  - **Access**: Admin Only
  - **Data** : Form Data
  - **Paramaters** : file

- **Search Manager Data**:

  - **Method**: POST
  - **Endpoint**: `http://localhost:8000/manager/search`
  - **Access**: Admin Only
  - **Search by**: First Name, Last Name, Gender
  - **Pagination**: `recordPerPage` (default: 10), `page` (default: 1)

- **Filter Managers**:
  - **Method**: POST
  - **Endpoint**: `http://localhost:8000/manager/filter`
  - **Access**: Admin Only
  - **Filter by**: age, salary, status
  - **Pagination**: `recordPerPage` (default: 10), `page` (default: 1)

### Guest Routes

- **Guest Signup**:

  - **Method**: POST
  - **Endpoint**: `http://localhost:8000/guest/signup`
  - **Access**: Public
  - **Data** : Form Data
  - **Paramaters** : first_name, last_name, phone, email, password, country, gender, id_proof_type

- **Guest Login**:

  - **Method**: POST
  - **Endpoint**: `http://localhost:8000/guest/login`
  - **Access**: Public
  - **Data** : Form Data
  - **Paramaters** : email and password

- **Verify Guest Email**:

  - **Method**: GET
  - **Endpoint**: `http://localhost:8000/guest/verify-email/confirm`
  - **Access**: Public (This route is email verification link that provided to email while register the user)

- **Get Guest Details**:

  - **Method**: GET
  - **Endpoint**: `http://localhost:8000/guest/get/:id`
  - **Access**: Guest Only

- **Update Guest Details**:

  - **Method**: PUT
  - **Endpoint**: `http://localhost:8000/guest/update/:id`
  - **Access**: Guest Only
  - **Data** : Form Data
  - **Paramaters** : first_name, last_name, phone, country, gender

- **Reset Guest Password**:

  - **Method**: PATCH
  - **Endpoint**: `http://localhost:8000/guest/update-password`
  - **Access**: Guest Only
  - **Data** : Form Data
  - **Paramaters** : email and password

- **Update Guest Profile Picture**:

  - **Method**: PATCH
  - **Endpoint**: `http://localhost:8000/guest/update-profile-pic/:id`
  - **Access**: Guest Only
  - **Data** : Form Data
  - **Paramaters** : file

- **Delete Guest Account**:

  - **Method**: DELETE
  - **Endpoint**: `http://localhost:8000/guest/delete/:id`
  - **Access**: Guest Only

- **Get All Guests**:
  - **Method**: GET
  - **Endpoint**: `http://localhost:8000/guest/getall`
  - **Access**: Admin Only
  - **Pagination**: `recordPerPage` (default: 10), `page` (default: 1)

### Driver Routes

- **Driver Login**:

  - **Method**: POST
  - **Endpoint**: `http://localhost:8000/driver/login`
  - **Access**: Public
  - **Data** : Form Data
  - **Paramaters** : email and password

- **Get all Drivers**:

  - **Method**: GET
  - **Endpoint**: `http://localhost:8000/driver/getall`
  - **Access**: Manager, Admin, Guest
  - **Pagination**: `recordPerPage` (default: 10), `page` (default: 1)

- **Get Driver**:

  - **Method**: GET
  - **Endpoint**: `http://localhost:8000/driver/get/:id`
  - **Access**: Manager, Admin, Guest
  - **Parameters**: driver_id

- **Create Driver**:

  - **Method**: POST
  - **Endpoint**: `http://localhost:8000/driver/create`
  - **Access**: Admin Only
  - **Data** : Form Data
  - **Paramaters** : first_name, last_name, email, password, gender, age, car_company, car_model, car_number_plate, phone, salary, file

- **Update Driver Details**:

  - **Method**: PUT
  - **Endpoint**: `http://localhost:8000/driver/update-all/:id`
  - **Access**: Admin Only
  - **Data** : Form Data
  - **Paramaters** : first_name, last_name, gender, age, car_company, car_model, car_number_plate, phone, salary

- **Update Driver Status**:

  - **Method**: PATCH
  - **Endpoint**: `http://localhost:8000/driver/update-status/:id`
  - **Access**: Admin Only

- **Delete Driver**:

  - **Method**: DELETE
  - **Endpoint**: `http://localhost:8000/driver/delete/:id`
  - **Access**: Admin Only

- **Update Driver Profile Picture**:

  - **Method**: PATCH
  - **Endpoint**: `http://localhost:8000/driver/update-profile-pic/:id`
  - **Access**: Driver Only
  - **Data** : Form Data
  - **Paramaters** : file

- **Reset Driver Password**:

  - **Method**: PATCH
  - **Endpoint**: `http://localhost:8000/driver/update-password`
  - **Access**: Driver Only
  - **Data** : Form Data
  - **Paramaters** : email and password

- **Update Driver Availability**:

  - **Method**: PATCH
  - **Endpoint**: `http://localhost:8000/driver/update-availability/:id`
  - **Access**: Driver Only
  - **Note**: Status is changed based on reservation

- **Search Driver Data**:

  - **Method**: POST
  - **Endpoint**: `http://localhost:8000/driver/search`
  - **Access**: Manager, Admin
  - **Search by**: First Name, Last Name, Gender, Car Company, Car Model, Car Number Plate
  - **Pagination**: `recordPerPage` (default: 10), `page` (default: 1)

- **Filter Driver**:
  - **Method**: POST
  - **Endpoint**: `http://localhost:8000/driver/filter`
  - **Access**: Manager, Admin
  - **Filter by**: availability, state, salary
  - **Pagination**: `recordPerPage` (default: 10), `page` (default: 1)

### Staff Routes

- **Get all Staff (By Branch ID)**:

  - **Method**: GET
  - **Endpoint**: `http://localhost:8000/staff/getall/:id`
  - **Access**: Manager
  - **Parameters**: branch_id
  - **Pagination**: `recordPerPage` (default: 10), `page` (default: 1)

- **Get all Staff**:

  - **Method**: GET
  - **Endpoint**: `http://localhost:8000/staff/getall`
  - **Access**: Manager
  - **Pagination**: `recordPerPage` (default: 10), `page` (default: 1)

- **Get Staff**:

  - **Method**: GET
  - **Endpoint**: `http://localhost:8000/staff/get/:id`
  - **Access**: Manager
  - **Parameters**: staff_id

- **Create Staff**:

  - **Method**: POST
  - **Endpoint**: `http://localhost:8000/staff/create`
  - **Access**: Manager
  - **Data** : Form Data
  - **Paramaters** : branch_id, first_name, last_name, phone, email, gender, age, job_type, salary, aadhar_number, file

- **Update Staff Details**:

  - **Method**: PUT
  - **Endpoint**: `http://localhost:8000/staff/update-all/:id`
  - **Access**: Manager
  - **Data** : Form Data
  - **Paramaters** : branch_id, first_name, last_name, phone, email, gender, age, job_type, salary, aadhar_number

- **Update Staff Status**:

  - **Method**: PATCH
  - **Endpoint**: `http://localhost:8000/staff/update-status/:id`
  - **Access**: Manager

- **Update Staff Profile Picture**:

  - **Method**: PATCH
  - **Endpoint**: `http://localhost:8000/staff/update-profile-pic/:id`
  - **Access**: Manager
  - **Data** : Form Data
  - **Paramaters** : file

- **Delete Staff**:

  - **Method**: DELETE
  - **Endpoint**: `http://localhost:8000/staff/delete/:id`
  - **Access**: Manager

- **Search Staff Data**:

  - **Method**: POST
  - **Endpoint**: `http://localhost:8000/staff/search`
  - **Access**: Manager
  - **Search by**: First Name, Last Name, Gender
  - **Pagination**: `recordPerPage` (default: 10), `page` (default: 1)

- **Filter Staff**:
  - **Method**: POST
  - **Endpoint**: `http://localhost:8000/staff/filter`
  - **Access**: Manager
  - **Search by**: age, salary, status, job_type
  - **Pagination**: `recordPerPage` (default: 10), `page` (default: 1)

### Feedback Routes

- **Get all Feedbacks (By Branch ID)**:

  - **Method**: GET
  - **Endpoint**: `http://localhost:8000/feedback/getall/:id`
  - **Access**: Manager, Admin
  - **Pagination**: `recordPerPage` (default: 10), `page` (default: 1)

- **Get Feedback**:

  - **Method**: GET
  - **Endpoint**: `http://localhost:8000/feedback/get/:id`
  - **Access**: Manager, Admin

- **Create Feedback**:

  - **Method**: POST
  - **Endpoint**: `http://localhost:8000/feedback/create`
  - **Access**: Manager, Guest
  - **Data** : Form Data
  - **Paramaters** : branch_id, room_id, guest_id, description, feedback_type, rating, file

- **Update Resolution Details**:

  - **Method**: PATCH
  - **Endpoint**: `http://localhost:8000/feedback/update-resolution-details/:id`
  - **Access**: Manager, Admin
  - **Data** : Form Data
  - **Paramaters** : resolution_details
  - **Note**: Only managers and admins can reply to feedback.

- **Delete Feedback**:

  - **Method**: DELETE
  - **Endpoint**: `http://localhost:8000/feedback/delete/:id`
  - **Access**: Guest
  - **Note**: Only guests can delete their own feedback.

- **Filter Feedback**:
  - **Method**: POST
  - **Endpoint**: `http://localhost:8000/feedback/filter`
  - **Access**: Manager, Admin
  - **Filter by**: status, feedback_type(rating, complaint), rating
  - **Pagination**: `recordPerPage` (default: 10), `page` (default: 1)

### Pickup Service Routes

- **Get all Pickup Services**:

  - **Method**: GET
  - **Endpoint**: `http://localhost:8000/service/getall`
  - **Access**: Manager, User, Driver
  - **Pagination**: `recordPerPage` (default: 10), `page` (default: 1)

- **Get Pickup Service**:

  - **Method**: GET
  - **Endpoint**: `http://localhost:8000/service/get/:id`
  - **Access**: Manager, User, Driver

- **Create Pickup Service**:

  - **Method**: POST
  - **Endpoint**: `http://localhost:8000/service/create`
  - **Access**: Manager, User, Driver
  - **Data** : Form Data
  - **Paramaters** : guest_id, branch_id, driver_id, pickup_location, pickup_time

- **Update Pickup Service Details**:

  - **Method**: PUT
  - **Endpoint**: `http://localhost:8000/service/update-details/:id`
  - **Access**: Manager, User, Driver
  - **Data** : Form Data
  - **Paramaters** : pickup_location, pickup_time

- **Update Pickup Service Status**:

  - **Method**: PATCH
  - **Endpoint**: `http://localhost:8000/service/update-status/:id`
  - **Access**: Manager, User, Driver
  - **Note**: Update status of a service (completed or not)

- **Delete Pickup Service**:
  - **Method**: DELETE
  - **Endpoint**: `http://localhost:8000/service/delete/:id`
  - **Access**: Manager, User, Driver
  - **Note**: Cancel the pickup service

### Reservation Routes

- **Get all Reservations**:

  - **Method**: GET
  - **Endpoint**: `http://localhost:8000/reservation/getall`
  - **Access**: Manager
  - **Pagination**: `recordPerPage` (default: 10), `page` (default: 1)

- **Get Reservation**:

  - **Method**: GET
  - **Endpoint**: `http://localhost:8000/reservation/get/:id`
  - **Access**: Manager
  - **Parameters**: reservation_id

- **Create Reservation**:

  - **Method**: POST
  - **Endpoint**: `http://localhost:8000/reservation/create`
  - **Access**: Manager, User
  - **Data** : Form Data
  - **Paramaters** : room_id, guest_id, check_in_time, check_out_time, desposit_amount, numbers_of_guests

- **Update Reservation Details**:

  - **Method**: PUT
  - **Endpoint**: `http://localhost:8000/reservation/update-all/:id`
  - **Access**: Manager, User
  - **Data** : Form Data
  - **Paramaters** : check_in_time, check_out_time, desposit_amount, numbers_of_guests, is_check_out
  - **Note**: Give details update otherwise other things remain as it is

- **Delete Reservation**:
  - **Method**: DELETE
  - **Endpoint**: `http://localhost:8000/reservation/delete/:id/:room_id`
  - **Access**: Manager, User
  - **Note**: Cancel the reservation

## Contributing

- Contributions are welcome! If you have any suggestions, bug reports, or feature requests, please open an issue or submit a pull request.

## Feedback

- If you have any feedback, please reach out to me at manankoyawala.dev@gmail.com

## Authors

- [@mananKoyawala](https://github.com/mananKoyawala)

## License

[MIT License](LICENSE)
