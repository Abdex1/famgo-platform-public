import {
  Injectable,
  NestInterceptor,
  ExecutionContext,
  CallHandler,
  Logger,
} from '@nestjs/common';
import { Observable } from 'rxjs';
import { tap } from 'rxjs/operators';
import { Request, Response } from 'express';

/**
 * Logging Interceptor
 * Logs all HTTP requests and responses with timing
 * FROM: DriveMind observability patterns
 */
@Injectable()
export class LoggingInterceptor implements NestInterceptor {
  private readonly logger = new Logger(LoggingInterceptor.name);

  intercept(context: ExecutionContext, next: CallHandler): Observable<any> {
    const request = context.switchToHttp().getRequest<Request>();
    const response = context.switchToHttp().getResponse<Response>();

    const startTime = Date.now();
    const { method, url, headers } = request;

    // Extract user info if authenticated
    const userId = request['user']?.sub || 'anonymous';

    // Log request
    this.logger.debug(
      `[${method}] ${url} - User: ${userId} - Headers: ${JSON.stringify(headers)}`,
    );

    return next.handle().pipe(
      tap(
        (data) => {
          const duration = Date.now() - startTime;
          const statusCode = response.statusCode;

          this.logger.log(
            `[${method}] ${url} - ${statusCode} - ${duration}ms - User: ${userId}`,
          );

          return data;
        },
        (error) => {
          const duration = Date.now() - startTime;
          const statusCode = error.status || 500;

          this.logger.error(
            `[${method}] ${url} - ${statusCode} - ${duration}ms - User: ${userId} - Error: ${error.message}`,
          );
        },
      ),
    );
  }
}
