const dateFormatter = new Intl.DateTimeFormat("es-MX", {
    day: "2-digit",
    month: "short",
    year: "numeric"
})

export function formatDateToShortReadable(date: Date) {
    return dateFormatter.format(date)
}
