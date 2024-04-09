import { SimCard } from "@mui/icons-material";
import {
  Box,
  Button,
  Card,
  Container,
  Divider,
  Stack,
  Tab,
  Tabs,
  Typography,
} from "@mui/material";
import NotificationItem from "./notification-item";

export default function Page() {
  return (
    <Container maxWidth={"xl"}>
      <Box
        sx={{ width: 500 }}
        pt={3}
        borderRight={"1px solid #00000008"}
        height={"100vh"}
      >
        <Stack
          direction={"row"}
          px={2}
          py={1}
          justifyContent={"space-between"}
          alignItems={"center"}
        >
          <Typography variant="subtitle2">Notification</Typography>
          <Typography
            component={Button}
            size="small"
            fontSize={14}
            px={1}
            color="primary"
          >
            Mark all as read
          </Typography>
        </Stack>
        <Box px={2}>
          <Tabs>
            <Tab label="All" />
            <Tab label="Danger" />
            <Tab label="Warning" />
            <Tab label="Health Check" />
          </Tabs>
          <Divider />
        </Box>
        <Stack sx={{ p: 2 }}>
          <NotificationItem />
        </Stack>
      </Box>
    </Container>
  );
}
