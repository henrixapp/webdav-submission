// tslint:disable
/**
 * Submissions
 * This API specifies the submissions service, as accessed by the web admin UI used by students, lecturers and tutors. 
 *
 * The version of the OpenAPI document: 1.0
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */

import { exists, mapValues } from '../runtime';
/**
 * a invitation to a submission
 * @export
 * @interface Invitation
 */
export interface Invitation  {
    /**
     * 
     * @type {string}
     * @memberof Invitation
     */
    id?: string;
    /**
     * 
     * @type {Date}
     * @memberof Invitation
     */
    createdAt?: Date;
    /**
     * 
     * @type {Date}
     * @memberof Invitation
     */
    updatedAt?: Date;
    /**
     * the invitedUserID is invited by invitingUserID
     * @type {number}
     * @memberof Invitation
     */
    invitedUserID?: number;
    /**
     * the invitingUserID invites an invitedUserID
     * @type {number}
     * @memberof Invitation
     */
    invitingUserID?: number;
    /**
     * the submission that is subject of the invitation
     * @type {string}
     * @memberof Invitation
     */
    submissionID?: string;
}

export function InvitationFromJSON(json: any): Invitation {
    return {
        'id': !exists(json, 'id') ? undefined : json['id'],
        'createdAt': !exists(json, 'createdAt') ? undefined : new Date(json['createdAt']),
        'updatedAt': !exists(json, 'updatedAt') ? undefined : new Date(json['updatedAt']),
        'invitedUserID': !exists(json, 'invitedUserID') ? undefined : json['invitedUserID'],
        'invitingUserID': !exists(json, 'invitingUserID') ? undefined : json['invitingUserID'],
        'submissionID': !exists(json, 'submissionID') ? undefined : json['submissionID'],
    };
}

export function InvitationToJSON(value?: Invitation): any {
    if (value === undefined) {
        return undefined;
    }
    return {
        'id': value.id,
        'createdAt': value.createdAt === undefined ? undefined : value.createdAt.toISOString(),
        'updatedAt': value.updatedAt === undefined ? undefined : value.updatedAt.toISOString(),
        'invitedUserID': value.invitedUserID,
        'invitingUserID': value.invitingUserID,
        'submissionID': value.submissionID,
    };
}


