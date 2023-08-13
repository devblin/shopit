import { useState } from 'react';
import { AppBar, Box, Toolbar, IconButton, Typography, MenuItem, Menu } from '@mui/material';
import { Add, MoreVert, Home } from '@mui/icons-material';

export default function NavBar() {
    const [mobileMoreAnchorEl, setMobileMoreAnchorEl] = useState(null);
    const menuItemAddBtnProps = { component: "a", href: "/add-item" };
    const menuItemHomeBtnProps = { component: "a", href: "/" };
    const isMobileMenuOpen = Boolean(mobileMoreAnchorEl);

    const handleMobileMenuClose = () => {
        setMobileMoreAnchorEl(null);
    };
    const handleMobileMenuOpen = (event) => {
        setMobileMoreAnchorEl(event.currentTarget);
    };

    const mobileMenuId = 'search-item-mobile';
    const renderMobileMenu = (
        <Menu
            anchorEl={mobileMoreAnchorEl}
            anchorOrigin={{
                vertical: 'top',
                horizontal: 'right',
            }}
            id={mobileMenuId}
            keepMounted
            transformOrigin={{
                vertical: 'top',
                horizontal: 'right',
            }}
            open={isMobileMenuOpen}
            onClose={handleMobileMenuClose}
        >
            <MenuItem {...menuItemHomeBtnProps}>
                <Home sx={{ marginRight: 1 }} />
                <p>Home</p>
            </MenuItem>
            <MenuItem {...menuItemAddBtnProps}>
                <Add sx={{ marginRight: 1 }} />
                <p>Add Item</p>
            </MenuItem>
        </Menu>
    );

    return (
        <Box sx={{ flexGrow: 1 }}>
            <AppBar position='fixed'>
                <Toolbar>
                    <Typography
                        component="div"
                        sx={{ display: 'flex', alignItems: 'end' }}
                    >
                        <img height={40} src={`${window.origin}/light-logo64.png`} alt="light-logo"></img>
                        <Typography variant='h5' marginLeft={1}>Inventory</Typography>
                    </Typography>
                    <Box sx={{ flexGrow: 1 }} />
                    <Box sx={{ display: { xs: 'none', md: 'flex' } }}>
                        <MenuItem {...menuItemHomeBtnProps}>
                            <Home sx={{ marginRight: 1 }} />
                            <p>Home</p>
                        </MenuItem>
                        <MenuItem {...menuItemAddBtnProps}>
                            <Add sx={{ marginRight: 1 }} />
                            <p>Add Item</p>
                        </MenuItem>
                    </Box>
                    <Box sx={{ display: { xs: 'flex', md: 'none' } }}>
                        <IconButton
                            size="large"
                            aria-label="show more"
                            aria-controls={mobileMenuId}
                            aria-haspopup="true"
                            onClick={handleMobileMenuOpen}
                            color="inherit"
                        >
                            <MoreVert />
                        </IconButton>
                    </Box>
                </Toolbar>
            </AppBar>
            {renderMobileMenu}
        </Box>
    );
}