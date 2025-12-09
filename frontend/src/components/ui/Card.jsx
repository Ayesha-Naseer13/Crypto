import { Card as MuiCard, CardContent, CardHeader, CardActions } from "@mui/material"

export default function Card({ title, subtitle, children, actions, sx = {}, ...props }) {
  return (
    <MuiCard sx={{ ...sx }} {...props}>
      {title && <CardHeader title={title} subheader={subtitle} />}
      <CardContent>{children}</CardContent>
      {actions && <CardActions>{actions}</CardActions>}
    </MuiCard>
  )
}
