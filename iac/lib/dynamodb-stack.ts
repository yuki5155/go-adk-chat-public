import * as cdk from 'aws-cdk-lib';
import * as dynamodb from 'aws-cdk-lib/aws-dynamodb';
import { Construct } from 'constructs';

/**
 * DynamoDB Stack Props
 */
export interface DynamoDBStackProps extends cdk.StackProps {
  projectName: string;
  environment: string;
}

/**
 * DynamoDB Stack for User Role Management and Chat System
 * Creates tables for user roles, role requests, and chat functionality
 */
export class DynamoDBStack extends cdk.Stack {
  public readonly userRolesTable: dynamodb.Table;
  public readonly roleRequestsTable: dynamodb.Table;
  public readonly chatThreadsTable: dynamodb.Table;
  public readonly chatSessionsTable: dynamodb.Table;
  public readonly chatEventsTable: dynamodb.Table;
  public readonly chatMemoriesTable: dynamodb.Table;

  constructor(scope: Construct, id: string, props: DynamoDBStackProps) {
    super(scope, id, props);

    const { projectName, environment } = props;

    // ============================================================================
    // User Roles Table
    // Stores user role assignments (excluding root user managed via env var)
    // ============================================================================
    this.userRolesTable = new dynamodb.Table(this, 'UserRolesTable', {
      tableName: `${projectName}-${environment}-user-roles`,
      partitionKey: {
        name: 'user_id',
        type: dynamodb.AttributeType.STRING
      },
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST, // On-demand billing for dev
      removalPolicy: environment === 'prod'
        ? cdk.RemovalPolicy.RETAIN
        : cdk.RemovalPolicy.DESTROY,
      pointInTimeRecovery: environment === 'prod', // Enable PITR for production
      encryption: dynamodb.TableEncryption.AWS_MANAGED,
      stream: dynamodb.StreamViewType.NEW_AND_OLD_IMAGES, // For audit logging
    });

    // GSI for email lookup
    this.userRolesTable.addGlobalSecondaryIndex({
      indexName: 'email-index',
      partitionKey: {
        name: 'email',
        type: dynamodb.AttributeType.STRING
      },
      projectionType: dynamodb.ProjectionType.ALL
    });

    // GSI for role lookup
    this.userRolesTable.addGlobalSecondaryIndex({
      indexName: 'role-index',
      partitionKey: {
        name: 'role',
        type: dynamodb.AttributeType.STRING
      },
      projectionType: dynamodb.ProjectionType.ALL
    });

    // ============================================================================
    // Role Requests Table
    // Tracks user requests for role subscriptions
    // ============================================================================
    this.roleRequestsTable = new dynamodb.Table(this, 'RoleRequestsTable', {
      tableName: `${projectName}-${environment}-role-requests`,
      partitionKey: {
        name: 'request_id',
        type: dynamodb.AttributeType.STRING
      },
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
      removalPolicy: environment === 'prod'
        ? cdk.RemovalPolicy.RETAIN
        : cdk.RemovalPolicy.DESTROY,
      pointInTimeRecovery: environment === 'prod',
      encryption: dynamodb.TableEncryption.AWS_MANAGED,
      stream: dynamodb.StreamViewType.NEW_AND_OLD_IMAGES,
    });

    // GSI for user_id lookup
    this.roleRequestsTable.addGlobalSecondaryIndex({
      indexName: 'user_id-index',
      partitionKey: {
        name: 'user_id',
        type: dynamodb.AttributeType.STRING
      },
      projectionType: dynamodb.ProjectionType.ALL
    });

    // GSI for status-based queries with timestamp sorting
    this.roleRequestsTable.addGlobalSecondaryIndex({
      indexName: 'status-requested_at-index',
      partitionKey: {
        name: 'status',
        type: dynamodb.AttributeType.STRING
      },
      sortKey: {
        name: 'requested_at',
        type: dynamodb.AttributeType.NUMBER
      },
      projectionType: dynamodb.ProjectionType.ALL
    });

    // ============================================================================
    // Chat Threads Table
    // Stores conversation threads for users
    // ============================================================================
    this.chatThreadsTable = new dynamodb.Table(this, 'ChatThreadsTable', {
      tableName: `${projectName}-${environment}-chat-threads`,
      partitionKey: {
        name: 'user_id',
        type: dynamodb.AttributeType.STRING
      },
      sortKey: {
        name: 'thread_id',
        type: dynamodb.AttributeType.STRING
      },
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
      removalPolicy: environment === 'prod'
        ? cdk.RemovalPolicy.RETAIN
        : cdk.RemovalPolicy.DESTROY,
      pointInTimeRecovery: environment === 'prod',
      encryption: dynamodb.TableEncryption.AWS_MANAGED,
    });

    // GSI for listing threads by updated_at
    this.chatThreadsTable.addGlobalSecondaryIndex({
      indexName: 'thread-updated-index',
      partitionKey: {
        name: 'user_id',
        type: dynamodb.AttributeType.STRING
      },
      sortKey: {
        name: 'updated_at',
        type: dynamodb.AttributeType.NUMBER
      },
      projectionType: dynamodb.ProjectionType.ALL
    });

    // ============================================================================
    // Chat Sessions Table
    // Stores ADK sessions for conversation context
    // ============================================================================
    this.chatSessionsTable = new dynamodb.Table(this, 'ChatSessionsTable', {
      tableName: `${projectName}-${environment}-chat-sessions`,
      partitionKey: {
        name: 'thread_id',
        type: dynamodb.AttributeType.STRING
      },
      sortKey: {
        name: 'session_id',
        type: dynamodb.AttributeType.STRING
      },
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
      removalPolicy: environment === 'prod'
        ? cdk.RemovalPolicy.RETAIN
        : cdk.RemovalPolicy.DESTROY,
      pointInTimeRecovery: environment === 'prod',
      encryption: dynamodb.TableEncryption.AWS_MANAGED,
    });

    // GSI for user sessions lookup
    this.chatSessionsTable.addGlobalSecondaryIndex({
      indexName: 'user-sessions-index',
      partitionKey: {
        name: 'user_id',
        type: dynamodb.AttributeType.STRING
      },
      sortKey: {
        name: 'created_at',
        type: dynamodb.AttributeType.NUMBER
      },
      projectionType: dynamodb.ProjectionType.ALL
    });

    // ============================================================================
    // Chat Events Table
    // Stores individual messages/events within sessions
    // ============================================================================
    this.chatEventsTable = new dynamodb.Table(this, 'ChatEventsTable', {
      tableName: `${projectName}-${environment}-chat-events`,
      partitionKey: {
        name: 'session_id',
        type: dynamodb.AttributeType.STRING
      },
      sortKey: {
        name: 'event_id',
        type: dynamodb.AttributeType.STRING
      },
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
      removalPolicy: environment === 'prod'
        ? cdk.RemovalPolicy.RETAIN
        : cdk.RemovalPolicy.DESTROY,
      pointInTimeRecovery: environment === 'prod',
      encryption: dynamodb.TableEncryption.AWS_MANAGED,
    });

    // ============================================================================
    // Chat Memories Table
    // Stores extracted memories for long-term context
    // ============================================================================
    this.chatMemoriesTable = new dynamodb.Table(this, 'ChatMemoriesTable', {
      tableName: `${projectName}-${environment}-chat-memories`,
      partitionKey: {
        name: 'thread_id',
        type: dynamodb.AttributeType.STRING
      },
      sortKey: {
        name: 'memory_id',
        type: dynamodb.AttributeType.STRING
      },
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
      removalPolicy: environment === 'prod'
        ? cdk.RemovalPolicy.RETAIN
        : cdk.RemovalPolicy.DESTROY,
      pointInTimeRecovery: environment === 'prod',
      encryption: dynamodb.TableEncryption.AWS_MANAGED,
    });

    // GSI for user memories lookup
    this.chatMemoriesTable.addGlobalSecondaryIndex({
      indexName: 'user-memories-index',
      partitionKey: {
        name: 'user_id',
        type: dynamodb.AttributeType.STRING
      },
      sortKey: {
        name: 'timestamp',
        type: dynamodb.AttributeType.NUMBER
      },
      projectionType: dynamodb.ProjectionType.ALL
    });

    // ============================================================================
    // CloudFormation Outputs
    // ============================================================================
    new cdk.CfnOutput(this, 'UserRolesTableName', {
      value: this.userRolesTable.tableName,
      description: 'User Roles Table Name',
      exportName: `${projectName}-${environment}-user-roles-table`
    });

    new cdk.CfnOutput(this, 'UserRolesTableArn', {
      value: this.userRolesTable.tableArn,
      description: 'User Roles Table ARN'
    });

    new cdk.CfnOutput(this, 'RoleRequestsTableName', {
      value: this.roleRequestsTable.tableName,
      description: 'Role Requests Table Name',
      exportName: `${projectName}-${environment}-role-requests-table`
    });

    new cdk.CfnOutput(this, 'RoleRequestsTableArn', {
      value: this.roleRequestsTable.tableArn,
      description: 'Role Requests Table ARN'
    });

    new cdk.CfnOutput(this, 'ChatThreadsTableName', {
      value: this.chatThreadsTable.tableName,
      description: 'Chat Threads Table Name',
      exportName: `${projectName}-${environment}-chat-threads-table`
    });

    new cdk.CfnOutput(this, 'ChatThreadsTableArn', {
      value: this.chatThreadsTable.tableArn,
      description: 'Chat Threads Table ARN'
    });

    new cdk.CfnOutput(this, 'ChatSessionsTableName', {
      value: this.chatSessionsTable.tableName,
      description: 'Chat Sessions Table Name',
      exportName: `${projectName}-${environment}-chat-sessions-table`
    });

    new cdk.CfnOutput(this, 'ChatSessionsTableArn', {
      value: this.chatSessionsTable.tableArn,
      description: 'Chat Sessions Table ARN'
    });

    new cdk.CfnOutput(this, 'ChatEventsTableName', {
      value: this.chatEventsTable.tableName,
      description: 'Chat Events Table Name',
      exportName: `${projectName}-${environment}-chat-events-table`
    });

    new cdk.CfnOutput(this, 'ChatEventsTableArn', {
      value: this.chatEventsTable.tableArn,
      description: 'Chat Events Table ARN'
    });

    new cdk.CfnOutput(this, 'ChatMemoriesTableName', {
      value: this.chatMemoriesTable.tableName,
      description: 'Chat Memories Table Name',
      exportName: `${projectName}-${environment}-chat-memories-table`
    });

    new cdk.CfnOutput(this, 'ChatMemoriesTableArn', {
      value: this.chatMemoriesTable.tableArn,
      description: 'Chat Memories Table ARN'
    });

    // Console output
    console.log('');
    console.log('=========================');
    console.log('DynamoDB Tables Created');
    console.log('=========================');
    console.log(`User Roles Table: ${this.userRolesTable.tableName}`);
    console.log('   - PK: user_id');
    console.log('   - GSI: email-index');
    console.log('   - GSI: role-index');
    console.log('');
    console.log(`Role Requests Table: ${this.roleRequestsTable.tableName}`);
    console.log('   - PK: request_id');
    console.log('   - GSI: user_id-index');
    console.log('   - GSI: status-requested_at-index');
    console.log('');
    console.log(`Chat Threads Table: ${this.chatThreadsTable.tableName}`);
    console.log('   - PK: user_id, SK: thread_id');
    console.log('   - GSI: thread-updated-index');
    console.log('');
    console.log(`Chat Sessions Table: ${this.chatSessionsTable.tableName}`);
    console.log('   - PK: thread_id, SK: session_id');
    console.log('   - GSI: user-sessions-index');
    console.log('');
    console.log(`Chat Events Table: ${this.chatEventsTable.tableName}`);
    console.log('   - PK: session_id, SK: event_id');
    console.log('');
    console.log(`Chat Memories Table: ${this.chatMemoriesTable.tableName}`);
    console.log('   - PK: thread_id, SK: memory_id');
    console.log('   - GSI: user-memories-index');
    console.log('=========================');
  }
}
