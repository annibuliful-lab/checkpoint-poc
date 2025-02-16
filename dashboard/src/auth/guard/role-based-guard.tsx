import { m } from "framer-motion";

import Container from "@mui/material/Container";
import Typography from "@mui/material/Typography";
import { Theme, SxProps } from "@mui/material/styles";

import { varBounce, MotionContainer } from "@/components/animate";

// ----------------------------------------------------------------------

type RoleBasedGuardProp = {
  hasContent?: boolean;
  roles?: string[];
  children: React.ReactNode;
  sx?: SxProps<Theme>;
};

export default function RoleBasedGuard({
  hasContent,
  roles,
  children,
  sx,
}: RoleBasedGuardProp) {
  // Logic here to get current user role
  const user = {
    role: "ADMIN",
  };

  // const currentRole = 'user';
  const currentRole = user?.role; // admin;

  if (typeof roles !== "undefined" && !roles.includes(currentRole)) {
    return hasContent ? (
      <Container
        component={MotionContainer}
        sx={{ textAlign: "center", ...sx }}
      >
        <m.div variants={varBounce().in}>
          <Typography variant="h3" sx={{ mb: 2 }}>
            Permission Denied
          </Typography>
        </m.div>

        <m.div variants={varBounce().in}>
          <Typography sx={{ color: "text.secondary" }}>
            You do not have permission to access this page
          </Typography>
        </m.div>
      </Container>
    ) : null;
  }

  return <> {children} </>;
}
