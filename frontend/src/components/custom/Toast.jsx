import { Snackbar, Alert } from '@mui/material';

export default function Toast(props) {
    return (
        <Snackbar
            anchorOrigin={{ vertical: "top", horizontal: "center" }}
            open={props.open}
            autoHideDuration={props.wait ? null : 1500}
            onClose={props.onClose}>
            <Alert
                variant="filled"
                severity={props.severity}
                sx={{ width: '100%' }}>
                {props.message}
            </Alert>
        </Snackbar>
    );
}