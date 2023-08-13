import Exception from "../custom/Exception"
import { SEVERITY } from "../../helpers/constants"

export default function NotFound(props) {
    return (
        <Exception
            severity={SEVERITY.ERROR}
            message={"Page Not Found"}
        >
        </Exception>
    )
}