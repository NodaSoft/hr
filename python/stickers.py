import requests
from PIL import Image
from io import BytesIO
import os
import math

def download_images(folder_name, destination_folder, oauth_token):
    # URL for accessing the folder on Yandex.Disk
    folder_url = f"https://cloud-api.yandex.net/v1/disk/resources/download?path={folder_name}"

    # Request headers with OAuth token
    headers = {"Authorization": f"OAuth {oauth_token}"}

    # Get the list of files in the folder
    response = requests.get(folder_url, headers=headers)
    if response.status_code != 200:
        print(f"Error retrieving the list of files from the folder {folder_name}: {response.text}")
        return

    files = response.json()

    # Create a folder to save images if it doesn't exist
    if not os.path.exists(destination_folder):
        os.makedirs(destination_folder)

    # Download images and save them
    for file_info in files:
        file_name = file_info["name"]
        file_url = file_info["file"]
        response = requests.get(file_url, headers=headers)
        if response.status_code != 200:
            print(f"Error downloading the file {file_name}: {response.text}")
            continue

        image_data = BytesIO(response.content)
        image = Image.open(image_data)

        # Save the image
        image.save(os.path.join(destination_folder, file_name))

    print(f"Images from the folder {folder_name} successfully downloaded to {destination_folder}")

def merge_images(image_folder, output_file):
    # Get the list of image files in the folder
    image_files = [f for f in os.listdir(image_folder) if f.endswith('.jpg') or f.endswith('.jpeg') or f.endswith('.png')]
    
    # Calculate the number of rows and columns for the grid
    num_images = len(image_files)
    num_cols = int(math.ceil(math.sqrt(num_images)))
    num_rows = int(math.ceil(num_images / num_cols))

    # Calculate the size of each image in the grid
    image_size = 200  # Adjust as needed
    grid_width = image_size * num_cols
    grid_height = image_size * num_rows

    # Create a new blank image for the grid
    tiff_image = Image.new('RGB', (grid_width, grid_height), color='white')

    # Merge images into the grid
    for i, image_file in enumerate(image_files):
        row = i // num_cols
        col = i % num_cols
        x_offset = col * image_size
        y_offset = row * image_size

        image = Image.open(os.path.join(image_folder, image_file))
        image.thumbnail((image_size, image_size))
        tiff_image.paste(image, (x_offset, y_offset))

    # Save the grid as a TIFF file
    tiff_image.save(output_file)

    print(f"Images successfully merged into the file {output_file}")


# Example usage
folder_name = "1388_12_Наклейки 3-D_3"
destination_folder = "images"
output_file = "Result.tif"
oauth_token = "your_oauth_token_here"

download_images(folder_name, destination_folder, oauth_token)
merge_images(destination_folder, output_file)
