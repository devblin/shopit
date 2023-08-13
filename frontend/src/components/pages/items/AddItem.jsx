import {
    Select,
    MenuItem,
    Typography,
    Stack,
    TextField,
    FormControl,
    Button,
    InputLabel,
    useMediaQuery,
} from '@mui/material';
import { useState } from 'react';
import {
    DEFAULT_IMAGE_URL,
    ITEMS_CATEGORY,
    SEVERITY,
    DEFAULT_IMAGE_NAME
} from '../../../helpers/constants';
import Toast from '../../custom/Toast';
import { addItem } from '../../../helpers/actions';

export default function AddItem() {
    const minWidth600 = useMediaQuery("(max-width:600px)");
    const imageActionProps = { fullWidth: true, variant: "contained" };
    const [alert, setAlert] = useState({
        open: false,
        severity: SEVERITY.SUCCESS,
        message: "",
        wait: false,
    });
    const [previewUrl, setPreviewUrl] = useState(DEFAULT_IMAGE_URL.SIZE400);
    const [itemDetails, setItemDetails] = useState({
        body: {
            Name: "",
            Price: 0,
            Stock: 0,
            Category: 0,
            Image: DEFAULT_IMAGE_NAME,
            Sold: 0
        }, image: {
            File: null
        }
    });
    const resetItemDetails = () => {
        setItemDetails({
            body: {
                Name: "",
                Price: 0,
                Stock: 0,
                Category: 0,
                Image: DEFAULT_IMAGE_NAME,
                Sold: 0,
            }, image: {
                File: null
            }
        });
        setPreviewUrl(DEFAULT_IMAGE_URL.SIZE400)
    };
    const handleItemChanage = (value, field) => {
        const tempItemDetails = Object.assign({}, itemDetails);
        tempItemDetails.body[field] = value;
        setItemDetails(tempItemDetails);
    }
    const handleImageChange = (e) => {
        const file = e.target.files[0];
        const reader = new FileReader();
        reader.readAsDataURL(file);
        const tempItemDetails = Object.assign({}, itemDetails);
        tempItemDetails.image.File = file;
        tempItemDetails.body.Image = file.name;

        reader.onloadend = (e) => {
            setItemDetails(tempItemDetails);
            setPreviewUrl(reader.result);
        }
    };
    const handleImageRemove = () => {
        const tempItemDetails = Object.assign({}, itemDetails);
        tempItemDetails.image.File = null;
        tempItemDetails.body.Image = DEFAULT_IMAGE_NAME;
        setPreviewUrl(DEFAULT_IMAGE_URL.SIZE400);
    }
    const renderItemsCategory = () => {
        const menuItem = [];

        for (const category in ITEMS_CATEGORY) {
            menuItem.push(<MenuItem
                key={category}
                value={category}
            >
                {ITEMS_CATEGORY[category]}
            </MenuItem>);
        };

        return menuItem;
    }
    const handleAlert = () => {
        setAlert({ open: false, severity: alert.severity, message: alert.message });
    }

    return (
        <Stack width={"100%"}>
            <Toast
                open={alert.open}
                wait={alert.wait}
                severity={alert.severity}
                message={alert.message}
                onClose={handleAlert}
            >
            </Toast>
            <Typography
                align='center'
                margin={3}
                variant="h4"
                color='text.secondary'
            >
                New Item
            </Typography>
            <Stack
                direction={minWidth600 ? "column" : "row"}
                alignItems={minWidth600 ? "center" : ""}
            >
                <Stack maxWidth={400} width="100%">
                    <img
                        height={400}
                        style={{ width: "100%", objectFit: "contain", border: "2px dotted", borderColor: "grey" }}
                        src={previewUrl}
                        alt="default-item"
                    />
                    <Stack direction={"row"} justifyContent="center">
                        <Button
                            {...imageActionProps}
                            color={SEVERITY.INFO}
                            sx={{ margin: 1 }}
                            component="label"
                        >
                            Choose
                            <input onChange={handleImageChange} type="file" accept="image/*" hidden />
                        </Button>
                        <Button
                            {...imageActionProps}
                            color={SEVERITY.ERROR}
                            sx={{ margin: 1 }}
                            onClick={handleImageRemove}
                        >
                            Remove
                        </Button>
                    </Stack>
                </Stack>
                <Stack width={"100%"} padding={1}>
                    <TextField
                        margin={"dense"}
                        label="Name"
                        value={itemDetails.body.Name}
                        onChange={(e) => {
                            handleItemChanage(e.target.value, "Name")
                        }}
                    ></TextField>
                    <TextField
                        type="number"
                        margin={"dense"}
                        label="Price"
                        value={itemDetails.body.Price}
                        onChange={(e) => {
                            handleItemChanage(e.target.valueAsNumber, "Price")
                        }}
                    ></TextField>
                    <TextField
                        type="number"
                        margin={"dense"}
                        label="Stock"
                        value={itemDetails.body.Stock}
                        onChange={(e) => {
                            handleItemChanage(e.target.valueAsNumber, "Stock")
                        }}
                    ></TextField>
                    <FormControl margin={"dense"}>
                        <InputLabel>Category</InputLabel>
                        <Select
                            label="Category"
                            value={itemDetails.body.Category}
                            onChange={(e) => {
                                handleItemChanage(parseInt(e.target.value), "Category")
                            }}
                        >
                            {renderItemsCategory()}
                        </Select>
                    </FormControl>
                    <Typography margin={1} color="text.secondary">
                        Note: Item image when uploaded is resized to as thumbnail of dimension 64x64 and standard image of dimension 400x400.
                    </Typography>
                    <Button
                        size='large'
                        variant='contained'
                        sx={{ padding: 2 }}
                        onClick={() => addItem(setAlert, itemDetails, resetItemDetails)}
                    >
                        Add
                    </Button>
                </Stack>
            </Stack>
        </Stack>
    );
}