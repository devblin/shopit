import { Grid, List, ListItem, Button, Avatar, ListItemAvatar, ListItemText, useMediaQuery } from '@mui/material';
import { useEffect, useState } from 'react';
import { Delete, Visibility, Edit } from '@mui/icons-material';
import { SEVERITY, MESSAGES } from '../../helpers/constants';
import Toast from '../custom/Toast';
import Exception from '../custom/Exception';
import { deleteItemById, getItemList } from '../../helpers/actions';
import { getImageBySize, getStockComponent } from '../../helpers/get';
import { usePageVisibility } from 'react-page-visibility';

export default function Items() {
    const minWidth600 = useMediaQuery("(max-width:600px)");
    const isVisible = usePageVisibility();
    const [items, setItems] = useState([]);
    const [exception, setException] = useState({ severity: SEVERITY.WARN, message: MESSAGES.INVENTORY_EMPTY });
    const [alert, setAlert] = useState({
        open: false,
        severity: SEVERITY.SUCCESS,
        message: ""
    });

    const handleAlert = () => {
        setAlert({ open: false, severity: alert.severity, message: alert.message });
    }

    const renderItems = () => {
        const itemCards = [];
        const buttonProps = { variant: "contained", edge: "end", sx: { marginLeft: 1 }, size: "small" };

        if (items && items.length > 0) {
            items.forEach((value) => {
                itemCards.push(
                    <ListItem
                        sx={{ padding: 0 }}
                        key={value.Id}
                        divider
                        secondaryAction={
                            [
                                <Button key={value.Id + "01"} href={`/view-item/${value.Id}`} {...buttonProps} color={SEVERITY.INFO}>
                                    <Visibility />
                                    {!minWidth600 && <code>View</code>}
                                </Button>,
                                <Button key={value.Id + "02"} href={`/edit-item/${value.Id}`} {...buttonProps} color={SEVERITY.WARN}>
                                    <Edit />
                                    {!minWidth600 && <code>Edit</code>}
                                </Button>,
                                <Button key={value.Id + "03"} onClick={() => deleteItemById(value.Id, setAlert, items, setItems)}   {...buttonProps} color={SEVERITY.ERROR}>
                                    <Delete />
                                    {!minWidth600 && <code>Delete</code>}
                                </Button>
                            ]
                        }
                    >
                        <ListItemAvatar key={value.Id + "0"}>
                            <Avatar>
                                <img height={64} src={getImageBySize(value.Image, 64)} alt='product'></img>
                            </Avatar>
                        </ListItemAvatar>
                        <ListItemText
                            key={value.Id + "1"}
                            primary={<div style={{ maxWidth: !minWidth600 ? "" : "20px" }}><b>{value.Name}</b> <code>${value.Price}</code></div>}
                            secondary={getStockComponent(value.Stock, "")}
                        />
                    </ListItem>
                )
            });
        } else {
            return <Exception severity={exception.severity} message={exception.message}></Exception>;
        }

        return itemCards;
    }

    useEffect(() => {
        getItemList(setAlert, setException, setItems);
    }, [isVisible]);

    return (
        <Grid item xs={12} width="100%">
            <Toast
                open={alert.open}
                severity={alert.severity}
                message={alert.message}
                onClose={handleAlert}
            ></Toast>
            <List>
                {renderItems()}
            </List>
        </Grid>
    );
}