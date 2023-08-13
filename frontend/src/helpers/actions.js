import axios from "axios";
import { API, MESSAGES, SEVERITY } from "./constants";
import { getStatusMessage, getImageBySize } from "./get";

let data, status, message = MESSAGES.SERVER_ERROR;


const getItemList = (setAlert, setException, setItems) => {
	axios.get(`${API.ITEMS}`)
		.then(res => {
			({ data, status } = res);
			({ message } = data);

			if (status === 200) {
				if (data.data && data.data.length > 0) {
					setItems(data.data);
				} else {
					setException({ severity: SEVERITY.WARN, message: MESSAGES.INVENTORY_EMPTY });
				}
			} else {
				setException({ severity: SEVERITY.ERROR, message: MESSAGES.ERROR_FETCHING });
				setAlert({ open: true, severity: SEVERITY.ERROR, message: data.message });
			}
		})
		.catch(e => {
			console.log(e);
			({ data, status } = e.response);
			({ message } = data);
			setAlert({ open: true, severity: SEVERITY.ERROR, message: getStatusMessage(status, message) });
		});
};

const deleteItemById = (itemId, setAlert, items = null, setItems = null) => {
	setAlert({ open: true, severity: SEVERITY.WARN, message: "Please wait...", wait: true });
	const body = { Id: itemId };

	axios.delete(API.ITEM, {
		data: body
	},)
		.then(res => {
			({ data, status } = res);
			({ message } = data);

			if (status === 200) {
				setAlert({ open: true, severity: SEVERITY.SUCCESS, message: getStatusMessage(status, message) });
				const tempItems = items.filter(item => item.Id !== itemId);
				setItems(tempItems);
			} else {
				setAlert({ open: true, severity: SEVERITY.ERROR, message: getStatusMessage(status, message) });
			}
		})
		.catch(e => {
			({ data, status } = e.response);
			({ message } = data);
			setAlert({ open: true, severity: SEVERITY.ERROR, message: getStatusMessage(status, message) });
		});
};

const uploadImage = (setAlert, imageFile, imageName, successMessage = "") => {
	const formData = new FormData();
	formData.append("image", imageFile, imageName);

	axios.post(`${API.ITEM_IMAGE}`, formData,)
		.then(res => {
			({ data, status } = res);
			({ message } = data);
			setAlert({ open: true, severity: SEVERITY.SUCCESS, message: getStatusMessage(status, successMessage || message) });
			return true;
		})
		.catch(e => {
			({ data, status } = e.response);
			({ message } = data);
			setAlert({ open: true, severity: SEVERITY.ERROR, message: getStatusMessage(status, message) });
			return false;
		});
};

const addItem = (setAlert, itemDetails, resetItemDetails) => {
	setAlert({ open: true, severity: SEVERITY.WARN, message: "Please wait...", wait: true });
	const { body } = itemDetails;
	const imageFile = itemDetails.image.File;
	let data, status, message = MESSAGES.SERVER_ERROR;

	axios.post(API.ITEM, body)
		.then(res => {
			({ data, status } = res);
			({ message } = data);

			if (status === 200) {
				if (imageFile != null) {
					uploadImage(setAlert, imageFile, data.data.Image, message);
					resetItemDetails();
				} else {
					setAlert({ open: true, severity: SEVERITY.SUCCESS, message: getStatusMessage(status, message) });
				}
			} else {
				setAlert({ open: true, severity: SEVERITY.ERROR, message: getStatusMessage(status, message) });
			}
		})
		.catch(e => {
			({ data, status } = e.response);
			({ message } = data);
			setAlert({ open: true, severity: SEVERITY.ERROR, message: getStatusMessage(status, message || MESSAGES.SERVER_ERROR) });
		});
}

const getItemDetails = async (itemId, setItemDetails, setAlert, setPreviewUrl = null) => {
	let data, status, message;

	axios.get(API.ITEM_DETAILS(itemId),)
		.then(res => {
			({ data, status } = res);
			({ data, message } = data);

			if (status === 200) {
				let itemDetails = { body: {}, image: {} };
				itemDetails.body.Category = data.Category;
				itemDetails.body.Image = data.Image;
				itemDetails.body.Price = data.Price;
				itemDetails.body.Stock = data.Stock;
				itemDetails.body.Name = data.Name;
				itemDetails.image.File = null;

				setItemDetails(itemDetails);
				if (setPreviewUrl != null) {
					setPreviewUrl(getImageBySize(data.Image, 400));
				}
			}
		})
		.catch(e => {
			({ data, status } = e.response);
			({ message } = data);
			setAlert({ open: true, severity: SEVERITY.ERROR, message: getStatusMessage(status, message || MESSAGES.SERVER_ERROR) });
		});
};

const updateItem = (itemId, itemDetails, setAlert) => {
	setAlert({ open: true, severity: SEVERITY.WARN, message: "Please wait...", wait: true });
	const { body } = itemDetails;
	body.Id = itemId;
	const { image } = itemDetails;
	let data, status, message = MESSAGES.SERVER_ERROR;

	axios.put(API.ITEM, body,)
		.then(res => {
			({ data, status } = res);
			({ message } = data);

			if (status === 200) {
				const imageFile = image.File;
				if (imageFile !== null) {
					const imageName = data.data.Image;
					uploadImage(setAlert, imageFile, imageName, message);
				} else {
					setAlert({ open: true, severity: SEVERITY.SUCCESS, message: getStatusMessage(status, message) });
				}
			} else {
				setAlert({ open: true, severity: SEVERITY.ERROR, message: getStatusMessage(status, message) });
			}
		})
		.catch(e => {
			({ data, status } = e.response);
			({ message } = data);
			setAlert({ open: true, severity: SEVERITY.ERROR, message: getStatusMessage(status, message || MESSAGES.SERVER_ERROR) });
		});
}

export {
	deleteItemById,
	getItemList,
	uploadImage,
	addItem,
	getItemDetails,
	updateItem
};