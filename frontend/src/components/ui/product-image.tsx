import { API_ORIGIN } from "@/services/api";
import { ImageIcon } from "lucide-react";

type ProductImageProps = {
  src?: string;
  name: string;
  className?: string;
};

export function ProductImage({ src, name, className = "h-14 w-14" }: ProductImageProps) {
  const imageSrc = resolveProductImage(src);

  if (imageSrc) {
    return (
      <img
        src={imageSrc}
        alt={name}
        className={`${className} shrink-0 rounded-md border border-line bg-white object-cover`}
        loading="lazy"
      />
    );
  }

  const initials = name
    .split(" ")
    .filter(Boolean)
    .slice(0, 2)
    .map((part) => part[0]?.toUpperCase())
    .join("");

  return (
    <div className={`${className} grid shrink-0 place-items-center rounded-md border border-line bg-brandSoft text-brand`}>
      {initials ? <span className="text-sm font-bold">{initials}</span> : <ImageIcon className="h-5 w-5" />}
    </div>
  );
}

function resolveProductImage(src?: string) {
  if (!src) {
    return "";
  }
  if (src.startsWith("http://") || src.startsWith("https://") || src.startsWith("data:") || src.startsWith("blob:")) {
    return src;
  }
  if (src.startsWith("/uploads/")) {
    return `${API_ORIGIN}${src}`;
  }
  return src;
}
