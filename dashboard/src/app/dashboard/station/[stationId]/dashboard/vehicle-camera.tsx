import Image from "@/components/image";
import {
  Announcement,
  AnnouncementOutlined,
  Contacts,
  ContactsOutlined,
  ReportProblem,
} from "@mui/icons-material";
import {
  Stack,
  Typography,
  Button,
  Divider,
  Box,
  Card,
  CardActions,
  CardContent,
} from "@mui/material";
import React from "react";

export default function VehicleCamera() {
  return (
    <Stack bgcolor={"#000"} height={1}>
      <Stack direction={"row"} height={1}>
        {[
          {
            label: "Driver image",
            placeholder: <ContactsOutlined sx={{ fontSize: 60 }} />,
          },
          {
            label: "License plate image",
            placeholder: <AnnouncementOutlined sx={{ fontSize: 60 }} />,
          },
        ].map((info) => (
          <Card key={info.label} sx={{ width: 1, height: 1 }}>
            <Stack height={1}>
              <Stack
                justifyContent={"center"}
                alignItems={"center"}
                flex={1}
                spacing={1}
              >
                {info.placeholder}
                <Typography variant="subtitle1" fontWeight={500}>
                  {info.label}
                </Typography>
              </Stack>
              <Stack alignItems={"end"} py={0.5} px={4} bgcolor={"#0000007A"}>
                <Typography variant="subtitle2" color={"#ffffff"}>
                  {info.label}
                </Typography>
              </Stack>
            </Stack>
          </Card>
        ))}
      </Stack>
      <Stack direction={"row"}>
        {[
          {
            label: "Front image",
            position: "FRONT",
          },
          {
            label: "Back image",
            position: "BACK",
          },
          {
            label: "Side image",
            position: "SIDE",
          },
        ].map((info) => (
          <Card key={info.label} sx={{ width: 1, height: 200 }}>
            <Stack height={1}>
              <Stack
                justifyContent={"center"}
                alignItems={"center"}
                flex={1}
                spacing={1}
              >
                <Stack
                  direction={info.position === "SIDE" ? "column" : "row"}
                  justifyContent={"end"}
                  alignItems={"end"}
                  spacing={0.5}
                >
                  {info.position === "BACK" && (
                    <Divider
                      sx={{ borderColor: "red", borderWidth: 3 }}
                      orientation="vertical"
                    />
                  )}
                  <Image src="/assets/icons/ic_vehicle.svg" alt="vehicle" />
                  {info.position === "FRONT" && (
                    <Divider
                      sx={{ borderColor: "red", borderWidth: 3 }}
                      orientation="vertical"
                    />
                  )}
                  {info.position === "SIDE" && (
                    <Divider
                      sx={{ borderColor: "red", borderWidth: 3, width: "50%" }}
                    />
                  )}
                </Stack>
                <Typography variant="subtitle1" fontWeight={500}>
                  {info.label}
                </Typography>
              </Stack>
              <Stack alignItems={"end"} py={0.5} px={4} bgcolor={"#0000007A"}>
                <Typography variant="subtitle2" color={"#ffffff"}>
                  {info.label}
                </Typography>
              </Stack>
            </Stack>
          </Card>
        ))}
      </Stack>
    </Stack>
  );
}
