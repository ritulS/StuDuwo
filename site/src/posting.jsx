import * as React from "react";
import house from "./house-unsplash.jpg";

async function formSubmit(e, url) {
  e.preventDefault();
  const form = document.querySelector("#listing_info");
  const formData = new FormData(form);
  try {
    const response = await fetch(`${url}/new_listing`, {
      method: "POST",
      body: formData,
    });
    window.location.reload();
  } catch (e) {
    console.error(e);
  }
}

export function Posting(props) {
  const [img, updateImage] = React.useState(house);
  const updateImagePreview = (e) => {
    const newImageSrc = URL.createObjectURL(e.target.files[0]);
    updateImage(newImageSrc);
  };

  const submit_listing_form = (e) => formSubmit(e, props.url);
  return (
    <div className="w-full max-w-full typo-bg min-h-screen">
      <div className="w-full flex justify-center pt-20">
        <div className="w-2/3 flex flex-row justify-center items-stretch">
          <div
            className="w-1/3 flex items-center bg-scroll"
            style={{
              backgroundImage: `url(${img})`,
              backgroundPosition: "center",
              backgroundSize: "cover",
            }}
          ></div>
          <form
            onSubmit={submit_listing_form}
            className="bg-white shadow-md rounded px-8 pt-6 pb-8"
            id="listing_info"
          >
            <div className="mb-4">
              <label
                className="block text-gray-700 text-sm font-bold mb-2"
                htmlFor="name"
              >
                Listing name:
              </label>
              <input
                className="shadow appearance-none border rounded w-full py-2"
                type="text"
                id="listing_name"
                name="listing_name"
                required={true}
              />
            </div>
            <div className="mb-4">
              <label
                className="block text-gray-700 text-sm font-bold mb-2"
                htmlFor="mail"
              >
                Email:
              </label>
              <input
                className="shadow appearance-none border rounded w-full py-2"
                type="email"
                id="owner_email"
                name="owner_email"
                required={true}
              />
            </div>
            <div className="mb-4">
              <label
                className="block text-gray-700 text-sm font-bold mb-2"
                htmlFor="add1"
              >
                Apartment no:
              </label>
              <input
                className="shadow appearance-none border rounded w-full py-2"
                type="text"
                id="address1"
                name="address1"
                required={true}
              />
            </div>
            <div className="mb-4">
              <label
                className="block text-gray-700 text-sm font-bold mb-2"
                htmlFor="add2"
              >
                Street:
              </label>
              <input
                className="shadow appearance-none border rounded w-full py-2"
                type="text"
                id="address2"
                name="address2"
                required={true}
              />
            </div>
            <div className="mb-4">
              <label
                className="block text-gray-700 text-sm font-bold mb-2"
                htmlFor="pin"
              >
                Pincode:
              </label>
              <input
                className="shadow appearance-none border rounded w-full py-2"
                type="text"
                id="pincode"
                name="pincode"
                required={true}
              />
            </div>
            <div className="mb-4">
              <label
                className="block text-gray-700 text-sm font-bold mb-2"
                htmlFor="rooms"
              >
                Rooms:
              </label>
              <input
                className="shadow appearance-none border rounded w-full py-2"
                type="text"
                id="rooms"
                name="rooms"
                required={true}
              />
            </div>
            <div className="mb-4">
              <label
                className="block text-gray-700 text-sm font-bold mb-2"
                htmlFor="price"
              >
                Price:
              </label>
              <input
                className="shadow appearance-none border rounded w-full py-2"
                type="text"
                id="price"
                name="price"
                required={true}
              />
            </div>
            <div className="mb-4">
              <label
                className="block text-gray-700 text-sm font-bold mb-2"
                htmlFor="img_id"
              >
                Listing img:
              </label>
              <input
                className="shadow appearance-none border rounded w-full py-2"
                type="file"
                id="img_id"
                name="img_id"
                required={true}
                onChange={updateImagePreview}
              />
            </div>
            <div>
              <button
                type="submit"
                className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-3 text-gray-700 mb-3 leading-tight focus:shadow-outline"
              >
                Upload listing
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
}
