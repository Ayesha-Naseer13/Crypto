import { TextField } from "@mui/material"

export default function Input({ label, error, helperText, ...props }) {
  return (
    <TextField
      label={label}
      error={Boolean(error)}
      helperText={error || helperText}
      fullWidth
      variant="outlined"
      {...props}
    />
  )
}
