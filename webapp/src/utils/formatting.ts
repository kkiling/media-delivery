export function formatDate(dateString?: string) {
  if (!dateString) return '';
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  });
}

export function getRatingColor(rating?: number) {
  if (!rating) return 'secondary';
  if (rating >= 7.8) return 'success';
  if (rating >= 5) return 'warning';
  return 'danger';
}
