# Shopit

A inventory management system, with backend written in **Golang**, with NoSQL database **MongoDB** and UI in React.js.

## Features:

As per challenge following features has been implemented.

-   Basic CRUD funtionalities:
    -   Create inventory items
    -   Edit Them
    -   Delete Them
    -   View a list of them
-   Allow image uploads and storing image with generated thumbnails.

## Usage:

### Production:

-   To test deployed site visit **[here](https://sho-pit.herokuapp.com/)**.
-   To check deployed site status visit **[here](https://sho-pit.herokuapp.com/api/status)**.

### Development:

To test the site locally, follow below steps:

> Note: Make sure you have Golang installed and configured in your system.

-   Clone the repo
-   Generate **`env/.env`** file using **`env/.env.example`** and do not change the existing env values.
-   Create MongoDB database with any name you want and assign value in **`.env`** file.
    ```
    DATABASE_URL=mongodb://127.0.0.1:27017
    DATABASE_NAME=
    PORT=5000
    AWS_ACCESS_KEY=
    AWS_SECRET_KEY=
    AWS_REGION=ap-south-1
    AWS_S3_BUCKET=shop-it
    GIN_MODE=debug
    DEFAULT_ITEM_IMAGE_NAME=default-product.jpg
    ENV=dev
    ```
-   Create AWS-S3 bucket with given **`AWS_S3_BUCKET=shop-it`** name, in **`AWS_REGION=ap-south-1`** region.
-   Obtain **`AWS_ACCESS_KEY`** and **`AWS_SECRET_KEY`** from your AWS account.
-   After configuring **`env/.env`** file, open your terminal and enter below command in root of this project dir.
    ```sh
    go run main.go
    ```

### API Endpoints:

Status API:

-   **GET /api/status** : To check API status.
    Sample Response:
    ```
    {
        "client":"171.60.156.251",
        "message":"Ok",
        "number_of_connections":0
    }
    ```

Item API:

-   **GET /api/item/list** : Get list of items in database.
    Sample Response:
    ```
    {
      "data": [
        {
          "_id": "61ec332bf410351028e0abc3",
          "category": 22,
          "image": "9dbda917-57a1-47e5-852a-0f64b091651a.jpeg",
          "name": "ITEM_NAME,
          "price": 2000,
          "sold": 0,
          "stock": 1
        }
      ],
      "description": "Items list",
      "message": "Success"
    }
    ```
-   **GET /api/item/:itemId** : Get details of item with item-id.
    Sample Response:
    ```
    {
      "data": {
        "_id": "61ec332bf410351028e0abc3",
        "category": 22,
        "image": "9dbda917-57a1-47e5-852a-0f64b091651a.jpeg",
        "name": "Football with pump",
        "price": 2000,
        "sold": 0,
        "stock": 1
      },
      "description": "Item details with given id",
      "message": "Success"
    }
    ```
-   **POST /api/item** : Add new item in database.
    Request Body:
    ```
    {
        "Category":1,
        "Name":"TestItem",
        "Price":4,
        "Stock":2,
        "Sold":0,
        "Image":"default-product.jpg"
    }
    ```
    200 - Response:
    ```
    {
      "data": {
        "InsertedID": "61ec45217bf302b0024fc3e3"
      },
      "description": "Add item",
      "message": "Item added successfully"
    }
    ```
-   **DELETE /api/item** : Delete item with item-id.
    Request Body:
    ```
    {
        "id":"61ec45217bf302b0024fc3e3"
    }
    ```
    200 - Response:
    ```
    {
      "data": {
        "DeletedCount": 1
      },
      "description": "Delete item with given id",
      "message": "Item successfully deleted"
    }
    ```
-   **PUT /api/item** : Update item details with item-id.
    Request Body:
    ```
    {
        "Category":1,
        "Name":"CHANGING_NAME",
        "Price":4,
        "Stock":2,
        "Sold":0,
        "Image":"default-product.jpg",
        "id":"61ec461f7bf302b0024fc3e4"
    }
    ```
    200 - Response:
    ```
    {
      "data": {
        "ImageID": "default-product.jpg",
        "MatchedCount": 1,
        "ModifiedCount": 1
      },
      "description": "Update item with given id",
      "message": "Item successfully updated"
    }
    ```
-   **POST /api/item/image** : Insert item image.
    To test this, use API testing tool such as [hoppscotch](https://hoppscotch.io/) or [thunder client](https://marketplace.visualstudio.com/items?itemName=rangav.vscode-thunder-client) to upload file, with file key/field name as **`image`**.
    Example (using [thunder client](https://marketplace.visualstudio.com/items?itemName=rangav.vscode-thunder-client)):
    ![](https://i.imgur.com/EKgP96S.png)
