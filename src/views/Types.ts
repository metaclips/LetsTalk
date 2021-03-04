declare global {
    export interface MediaDevices {
        getDisplayMedia(constraints?: MediaStreamConstraints): Promise<MediaStream>;
    }

    // if constraints config still lose some prop, you can define it by yourself also
    export interface MediaTrackConstraintSet {
        displaySurface?: ConstrainDOMString;
        logicalSurface?: ConstrainBoolean;
        // more....
    }
}

export interface RoomDetails {
    roomID: string;
    roomName: string;
}

export interface JoinedRoom extends RoomDetails {
    roomIcon: string;
    joined: boolean;
    userID: string;
    name: string;
}

export interface UsersDetails {
    names: string;
    email: string;
    department: string;
    entryYear: number;
    dob: string;
}

export interface JoinRequest extends RoomDetails {
    requestingUserName: string;
    requestingUserID: string;
}

export interface FetchedUsers {
    name: string;
    userID: string;
}



export interface SentRoomRequest extends RoomDetails {
    requesterID: string;
    requesterName: string;
    userRequested: string;
}

export interface NewRoomRequest extends RoomDetails {
    requestingUserName: string;
    requestingUserID: string;
}

export interface RoomPageDetails {
    roomID: string;
    roomName: string;
    roomIcon: string;

    firstLoad: boolean;

    registeredUsers: string[];
    messages: Message[];
}

export interface Message {
    time?: string;
    type?: string; // Type denotes the message type, if INFO, FILE or MESSAGE type.
    name: string;
    userID: string;
    roomID: string;
    message: string;
    size?: string;
    hash?: string;

    index: number;
}

export interface FileDownload {
    [index: string]: FileUploadDownloadDetails;
}

export interface FileUploadDownloadDetails {
    roomID: string;
    userID: string;
    fileName: string;
    fileType?: string;
    fileHash: string;
    // File hash of compressed file if files are sent in chunks
    compressedFileHash?: string;

    downloading: boolean;
    isDownloader: boolean;
    fileSize: string;

    progress: number;
    chunks: number;
    chunk: number;
    nextChunk: number;

    fileContent?: any;
    fileChunk?: string;
}

export interface FileUploadComplete {
    userID: string;
    name: string;
    fileName: string;
    fileSize: string;
    fileHash: string;
    roomID: string;
}

export interface ExitRoomDetails {
    roomID: string;
    userID: string;
}

export interface FileChunk {
    fileName: string;
}

export interface AssociateStatus {
    name: string;
    isOnline: boolean;
}

export interface UsersIcon {
    [index: string]: string;
}

export interface RecentChatPreview {
    [index: string]: string;
}

export interface UsersOnline {
    [index: string]: AssociateStatus;
}

export interface UnreadRooms {
    [index: string]: boolean;
}