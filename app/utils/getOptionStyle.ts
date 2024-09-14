export const getOptionStyle = (type: "user" | "group") => {
    return type === "user" ? { color: "blue" } : { color: "green" };
};
