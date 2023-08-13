import { Typography } from "@mui/material";
import { CheckCircle, Error, Warning, Info, } from "@mui/icons-material";
import { SEVERITY, IMAGE_BASE_URL } from "./constants";

const getItemImageUrl = (image) => {
	const ITEM_IMAGE_BASE_URL = process.env.REACT_APP_ITEM_IMAGE_BASE_URL;
	return ITEM_IMAGE_BASE_URL + image;
}

const getStockComponent = (value, size) => {
	if (value > 0) {
		return <Typography fontSize={size} color="green">In-Stock: {value}</Typography>;
	}
	return <Typography fontSize={size} color="red">Out-of-Stock: 0</Typography>
}

const getPriceComponent = (value) => {
	return <Typography variant='h4'>${value}</Typography>;
}

const getExceptionIcon = (value) => {
	const props = { fontSize: "large", color: value || "success" }
	let icon = <CheckCircle {...props}></CheckCircle>;

	if (value === SEVERITY.ERROR) {
		icon = <Error {...props}></Error>
	} else if (value === SEVERITY.SUCCESS) {
		icon = <CheckCircle  {...props}></CheckCircle>
	} else if (value === SEVERITY.INFO) {
		icon = <Info {...props}></Info>
	} else {
		icon = <Warning {...props}></Warning>
	}

	return icon;
}

const getFileExtension = (file) => {
	const fileNameSplit = file.name.split(".");
	const fileNameSplitLen = fileNameSplit.length;
	const fileExtension = fileNameSplit[fileNameSplitLen - 1];

	return fileExtension;
}

const getImageName = (name) => {
	let imageNameSplit = name.split(".");
	const imageNameSplitLen = imageNameSplit.length;
	imageNameSplit.splice(imageNameSplitLen - 1, 1);
	imageNameSplit = imageNameSplit.join("");

	return imageNameSplit;
}

const getNewImageName = (name, ext) => {
	return name + "." + ext;
}

const getStatusMessage = (status, message) => {
	return `${status} : ${message}`;
}

const getImageBySize = (name, size) => {
	let imageNameSplit = name.split(".");
	const imageNameSplitLen = imageNameSplit.length;
	const imageExtension = imageNameSplit[imageNameSplitLen - 1];
	imageNameSplit.splice(imageNameSplitLen - 1, 1);
	imageNameSplit = imageNameSplit.join("");
	const imageName = imageNameSplit;

	return `${IMAGE_BASE_URL}/${imageName}${size}.${imageExtension}`;
}

export {
	getItemImageUrl,
	getStockComponent,
	getPriceComponent,
	getExceptionIcon,
	getFileExtension,
	getNewImageName,
	getStatusMessage,
	getImageBySize,
	getImageName,
}