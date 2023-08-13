import { Stack, Typography, Link } from "@mui/material"

export default function FootBar() {
    return (
        <Stack padding={6} alignItems="center">
            <Typography color="text.secondary">
                Made by <Link href="https://github.com/devblin" target={"_blank"}>devblin</Link> | © 2022 Shopit
            </Typography>
            <Typography fontSize={"12px"} color="text.secondary">
                This site is for educational purposes only.
            </Typography>
        </Stack>
    )
}