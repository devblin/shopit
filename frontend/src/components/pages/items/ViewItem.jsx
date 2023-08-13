import { useEffect, useState } from "react";
import Item from "./Item";
import { ITEMS_CATEGORY, SEVERITY, DEFAULT_IMAGE_NAME } from "../../../helpers/constants";
import Toast from "../../custom/Toast";
import { getItemDetails } from "../../../helpers/actions";
import { usePageVisibility } from 'react-page-visibility';

export default function ViewItem(props) {
    const isVisible = usePageVisibility();
    const [itemId] = useState(props.params.itemId);
    const [alert, setAlert] = useState({
        open: false,
        severity: SEVERITY.SUCCESS,
        message: "",
        wait: false,
    });
    const [itemDetails, setItemDetails] = useState({
        body: {
            Name: "Item doesn't exists",
            Image: DEFAULT_IMAGE_NAME,
            Stock: 0,
            Price: 0,
            Category: 0,
        }
    });

    const handleAlert = () => {
        setAlert({ open: false, severity: alert.severity, message: alert.message });
    }

    useEffect(() => {
        getItemDetails(itemId, setItemDetails, setAlert);
    }, [isVisible, itemId]);

    return (
        <div>
            <Toast
                open={alert.open}
                severity={alert.severity}
                message={alert.message}
                onClose={handleAlert}
            ></Toast><Item
                name={itemDetails.body.Name}
                image={itemDetails.body.Image}
                category={ITEMS_CATEGORY[itemDetails.body.Category]}
                price={itemDetails.body.Price}
                stock={itemDetails.body.Stock}
                id={itemId}
            ></Item>
        </div>)
}