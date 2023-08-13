import { Stack, Typography } from "@mui/material"
import { getExceptionIcon } from "../../helpers/get"

export default function Exception(props) {
    return (
        <Stack padding={6} alignItems="center">
            <Typography variant="h5" textAlign={"center"} maxWidth="500px">
                {getExceptionIcon(props.severity)}
                <Typography marginTop={1} fontSize={25}>
                    {props.message}
                </Typography>
            </Typography>
        </Stack>
    )
}