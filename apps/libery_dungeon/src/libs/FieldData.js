export const FieldStates = {
    NORMAL: 0,
    HAS_ERRORS: 1,
    READY: 2
}

export const verifyFormFields = (form_data) => {
    let field_value = null;
    let is_valid = true;
    for(let fd of form_data) {
        if (fd instanceof Array) {
            is_valid = verifyFormFields(fd);
            continue;
        }

        field_value = fd.getFieldValue();
        if (field_value === "") {
            
            if (fd.required) {
                fd.state = FieldStates.NORMAL;
                is_valid = false;
            } else {
                fd.state = FieldStates.READY;
            }
        } else if(!fd.isReady()) {

            fd.state = FieldStates.HAS_ERRORS;
            fd.error_message = "Invalid field";
            is_valid = false;
        } else {
            fd.state = FieldStates.READY;
        }
    }

    return is_valid;
}

export const resetFormFields = (form_data) => {
    for(let fd of form_data) {
        fd.state = FieldStates.NORMAL;
        fd.setFieldValue("");
        fd.error_message = "";
    }
}

class FieldData {
    constructor(field_id, validation_regex, name,type_name="text", required=true) {
        this.id = field_id;
        this.name = name;   
        this.regex = validation_regex;
        this.type = type_name;
        this.state = FieldStates.NORMAL;
        this.required = required;
        this.error_message = "hey! im error"; // link's adventure...
        this.placeholder = this.name;
    }

    clear = () => {
        let field = this.getField();

        if (field !== null) {
            field.value = '';
        }
    }

    getField = () => {
        return document.getElementById(this.id);
    }

    getFieldValue = () => {
        let field = this.getField();
        
        if (!field) {
            return "";
        }

        return this.getField().value;
    }

    setFieldValue = (value) => {
        let field = this.getField();

        if (!field) {
            return;
        }

        field.value = value;
    }

    isReady = () => {
        return this.regex.test(this.getFieldValue());
    }

    setPlaceholder = placeholder => {
        this.placeholder = placeholder;
    }

    verify = () => {
        if (this.getFieldValue() === "") {
            if (this.required) {
                this.state = FieldStates.NORMAL;
                return false;
            } else {
                this.state = FieldStates.READY;
                return true;
            }
        } else if(!this.isReady()) {
            this.state = FieldStates.HAS_ERRORS;
            this.error_message = "Invalid field";
            return false;
        } else {
            this.state = FieldStates.READY;
            return true;
        }
    }
}

export default FieldData;