import { Alert } from 'react-bootstrap';

export default function NotFound() {
  return (
    <Alert variant="danger">
      <Alert.Heading>Page Not Found</Alert.Heading>
      <p>The page youre looking for doesnt exist.</p>
    </Alert>
  );
}
